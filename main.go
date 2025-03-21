package main

import (
	"flag"
	"os"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	database "github.com/mazay/mikromanager/db"
	"github.com/mazay/mikromanager/http"
	"github.com/mazay/mikromanager/internal"
	"go.uber.org/zap"
)

type PollerCFG struct {
	Client *internal.Api
	Db     *database.DB
	Device *database.Device
}

type BackupCFG struct {
	Client *internal.SshClient
	Db     *database.DB
	Device *database.Device
}

var (
	err        error
	configPath string
	httpPort   string

	s3 *internal.S3

	policy = &database.ExportsRetentionPolicy{Name: "Default"}

	logger = &zap.Logger{}
	wg     = sync.WaitGroup{}
	osExit = os.Exit
)

func main() {
	// Read command line args
	flag.StringVar(&configPath, "config", "config.yml", "Path to the config.yml")
	flag.StringVar(&httpPort, "http-port", "8080", "Port for the HTTP server")
	flag.Parse()

	config := readConfigFile(configPath)
	logger = initLogger(config.LogLevel)
	defer logger.Sync() //nolint:golint,errcheck

	pollerCH := make(chan *PollerCFG)
	exportCH := make(chan *BackupCFG)

	// init S3 client
	s3 = &internal.S3{
		Bucket:          config.S3Bucket,
		BucketPath:      config.S3BucketPath,
		Endpoint:        config.S3Endpoint,
		Region:          config.S3Region,
		StorageClass:    config.S3StorageClass,
		AccessKey:       config.S3AccessKey,
		SecretAccessKey: config.S3SecretAccessKey,
		OpsRetries:      config.S3OpsRetries,
	}
	err = s3.GetS3Session()
	if err != nil {
		logger.Panic("S3 client init issue", zap.String("error", err.Error()))
		osExit(1)
	}

	wg.Add(1)

	db := database.DB{LogLevel: config.DbLogLevel}
	err = db.Open(config.DbPath)
	if err != nil {
		logger.Panic("DB init issue", zap.String("error", err.Error()))
		osExit(1)
	}
	defer db.Close()

	logger.Debug("ensure at least one user exists, create 'admin' otherwise")
	user := &database.User{}
	users, err := user.GetAll(&db)
	if err != nil {
		logger.Error(err.Error())
	}

	if len(users) == 0 {
		logger.Info("no users found, creating 'admin' user")
		encryptedPw, err := database.EncryptString("admin", config.EncryptionKey)
		if err != nil {
			logger.Error(err.Error())
			osExit(3)
		}
		user.Username = "admin"
		user.EncryptedPassword = encryptedPw
		err = user.Create(&db)
		if err != nil {
			logger.Error(err.Error())
			osExit(3)
		}
	}

	// run HTTP server
	server := http.HttpConfig{
		Port:          "8000",
		Db:            &db,
		EncryptionKey: config.EncryptionKey,
		Logger:        logger,
		BackupPath:    config.BackupPath,
		S3:            s3,
	}
	go server.HttpServer()

	logger.Info("starting MikroTik API pollers", zap.Int("count", config.ApiPollers))
	for range make([]int, config.ApiPollers) {
		wg.Add(1)
		go apiWorker(pollerCH)
	}

	logger.Info("starting MikroManager export workers", zap.Int("count", config.ExportWorkers))
	for range make([]int, config.ExportWorkers) {
		wg.Add(1)
		go exportWorker(exportCH)
	}

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		logger.Error("scheduler", zap.Any("error", err))
	}
	logger.Info("devicePollerInterval", zap.Duration("interval", config.DevicePollerInterval))
	pollerJob, pollerErr := scheduler.NewJob(
		gocron.DurationJob(config.DevicePollerInterval),
		gocron.NewTask(devicesPoller, config, &db, pollerCH),
	)
	if pollerErr != nil {
		logger.Error("poller", zap.Any("Job", pollerJob), zap.Any("error", pollerErr))
	}
	logger.Info("deviceExportCronSchedule", zap.String("cron schedule", config.deviceExportCronSchedule))
	exportJob, exportErr := scheduler.NewJob(
		gocron.CronJob(config.deviceExportCronSchedule, false),
		gocron.NewTask(backupScheduler, config, &db, exportCH),
	)
	if exportErr != nil {
		logger.Error("export", zap.Any("Job", exportJob), zap.Any("error", exportErr))
	}
	logger.Info("export retention job interval is 24 hours")
	exportRetentionJob, exportRetentionErr := scheduler.NewJob(
		gocron.CronJob("0 * * * *", false),
		gocron.NewTask(rotateExports, &db),
	)
	if exportRetentionErr != nil {
		logger.Error("export", zap.Any("Job", exportRetentionJob), zap.Any("error", exportRetentionErr))
	}
	logger.Info("session cleanup job interval is 24 hours")
	sessionCleanupJob, sessionCleanupErr := scheduler.NewJob(
		gocron.CronJob("0 * * * *", false),
		gocron.NewTask(cleanupSessions, &db),
	)
	if sessionCleanupErr != nil {
		logger.Error("session", zap.Any("Job", sessionCleanupJob), zap.Any("error", sessionCleanupErr))
	}
	scheduler.Start()

	wg.Wait()
}

func devicesPoller(cfg *Config, db *database.DB, pollerCH chan<- *PollerCFG) error {
	var d = &database.Device{}

	logger.Info("starting device polling task")
	// have to use the preload method here to not lose the device groups
	devices, err := d.GetAllPlain(db)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	for _, device := range devices {
		creds, err := device.GetCredentials(db)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		logger.Debug("authentication", zap.String("credentials", creds.Alias), zap.String("device", device.Address))
		decryptedPw, encryptionErr := database.DecryptString(creds.EncryptedPassword, cfg.EncryptionKey)
		if encryptionErr != nil {
			logger.Error(encryptionErr.Error())
			return err
		}
		client := &internal.Api{
			Address:  device.Address,
			Port:     device.ApiPort,
			Username: creds.Username,
			Password: decryptedPw,
			Async:    true,
			UseTLS:   false,
			Logger:   logger,
		}
		pollerCH <- &PollerCFG{Client: client, Db: db, Device: device}
	}
	return nil
}

func apiWorker(pollerCH <-chan *PollerCFG) {
	for cfg := range pollerCH {
		var fetchErr error
		var minorErr error
		var dbErr error

		logger.Info("polling device", zap.String("address", cfg.Client.Address))
		fetchErr = fetchResources(cfg)
		if fetchErr != nil {
			logger.Error(fetchErr.Error())
		}

		fetchErr = fetchRbDetails(cfg)
		if fetchErr != nil {
			logger.Error(fetchErr.Error())
		}

		fetchErr = fetchIdentity(cfg)
		if fetchErr != nil {
			logger.Error(fetchErr.Error())
		}

		minorErr = cfg.Client.CheckForUpdates(cfg.Device)
		if minorErr != nil {
			logger.Error(minorErr.Error())
		}

		// do not consider fetchManagementIp errors as a failure, just log them
		minorErr = fetchManagementIp(cfg)
		if minorErr != nil {
			logger.Error(minorErr.Error())
		}

		if fetchErr != nil {
			cfg.Device.PollingSucceeded = 0
		} else {
			cfg.Device.PollingSucceeded = 1
			cfg.Device.PolledAt = time.Now()
		}

		dbErr = cfg.Device.Save(cfg.Db)

		if dbErr != nil {
			logger.Error(dbErr.Error())
		}
	}
}

func backupScheduler(cfg *Config, db *database.DB, exportCH chan<- *BackupCFG) {
	var d = &database.Device{}

	logger.Info("starting backup task")
	devices, err := d.GetAllPlain(db)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, device := range devices {
		creds, err := device.GetCredentials(db)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		logger.Debug("authentication", zap.String("credentials", creds.Alias), zap.String("device", device.Address))
		decryptedPw, encryptionErr := database.DecryptString(creds.EncryptedPassword, cfg.EncryptionKey)
		if encryptionErr != nil {
			logger.Error(encryptionErr.Error())
			return
		}
		client := &internal.SshClient{
			Host:     device.Address,
			Port:     device.SshPort,
			User:     creds.Username,
			Password: decryptedPw,
		}
		exportCH <- &BackupCFG{Client: client, Db: db, Device: device}
	}
}

func exportWorker(exportCH <-chan *BackupCFG) {
	for cfg := range exportCH {
		logger.Debug("creating backup", zap.String("address", cfg.Client.Host))

		export, sshErr := cfg.Client.Run("/export show-sensitive")
		if sshErr == nil {
			output, err := s3.UploadExport(cfg.Device.Id, []byte(export))
			if err != nil {
				logger.Error(err.Error())
			}

			attrs, err := s3.GetExportAttributes(*output.Key)
			if err != nil {
				logger.Error(err.Error())
			}

			export := &database.Export{
				S3Key:        *output.Key,
				LastModified: attrs.LastModified,
				ETag:         *output.ETag,
				Size:         attrs.Size,
				DeviceId:     cfg.Device.Id,
			}
			err = export.Save(cfg.Db)
			if err != nil {
				logger.Error(err.Error())
			}

			logger.Info("created a new backup", zap.String("device", cfg.Device.Address), zap.String("s3 key", *output.Key))
		} else {
			logger.Error(sshErr.Error())
		}
	}
}

func rotateExports(db *database.DB) {
	var (
		err         error
		export      *database.Export
		exportsList []*database.Export
		device      *database.Device
	)

	logger.Info("starting exports retention task")
	err = policy.GetDefault(db)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	devices, err := device.GetAllPreload(db)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, device := range devices {
		exports, err := export.GetByDeviceId(db, device.Id)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		exportsList = append(exportsList, rotateHourlyExports(exports, policy.Hourly)...)
		exportsList = append(exportsList, rotateDailyExports(exports, policy.Daily)...)
		exportsList = append(exportsList, rotateWeeklyExports(exports, policy.Weekly)...)

		for _, export := range exports {
			if !exportInSlice(export, exportsList) {
				logger.Debug("deleting export", zap.String("filename", export.S3Key))

				err := s3.DeleteFile(export.S3Key)
				if err != nil {
					logger.Error(err.Error())
				}

				err = export.Delete(db)
				if err != nil {
					logger.Error(err.Error())
				}
			}
		}
	}
}

func cleanupSessions(db *database.DB) {
	var err error
	var session *database.Session

	logger.Info("starting session cleanup task")
	sessions, err := session.GetAll(db)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	for _, s := range sessions {
		if s.ValidThrough.Before(time.Now()) {
			logger.Debug("session expired", zap.String("id", s.Id))
			err = s.Delete(db)
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}
}
