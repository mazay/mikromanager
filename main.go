package main

import (
	"flag"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/mazay/mikromanager/api"
	"github.com/mazay/mikromanager/db"
	"github.com/mazay/mikromanager/http"
	"github.com/mazay/mikromanager/internal"
	"github.com/mazay/mikromanager/ssh"
	"github.com/mazay/mikromanager/utils"
	"go.uber.org/zap"
)

type PollerCFG struct {
	Client *api.API
	Db     *db.DB
	Device *utils.Device
}

type BackupCFG struct {
	Client *ssh.SshClient
	Db     *db.DB
	Device *utils.Device
}

var (
	err        error
	configPath string
	httpPort   string

	s3 *internal.S3

	policy = &utils.ExportsRetentionPolicy{Name: "Default"}
	user   = &utils.User{}

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

	db := &db.DB{Path: config.DbPath}
	err = db.Init()
	if err != nil {
		logger.Panic("DB init issue", zap.String("error", err.Error()))
		osExit(1)
	}
	defer db.Close()

	logger.Debug("ensure 'Default' exports retention policy exists")
	err = policy.GetDefault(db)
	if err != nil || policy.Id == "" {
		logger.Error(err.Error())
	}

	if policy.Id == "" {
		policy.Hourly = 24
		policy.Daily = 14
		policy.Weekly = 26
		err = policy.Create(db)
		if err != nil {
			logger.Fatal(err.Error())
		}
	}

	logger.Debug("ensure at least one user exists, create 'admin' otherwise")
	users, err := user.GetAll(db)
	if err != nil {
		logger.Error(err.Error())
	}

	if len(users) == 0 {
		encryptedPw, err := utils.EncryptString("admin", config.EncryptionKey)
		if err != nil {
			logger.Error(err.Error())
			osExit(3)
		}
		user.Username = "admin"
		user.EncryptedPassword = encryptedPw
		err = user.Create(db)
		if err != nil {
			logger.Error(err.Error())
			osExit(3)
		}
	}

	collections, _ := db.ListCollections()
	logger.Debug("DB collections", zap.String("list", strings.Join(collections, ", ")))

	go http.HttpServer("8000", db, config.EncryptionKey, config.BackupPath, logger, s3)

	scheduler := gocron.NewScheduler(time.Local)
	logger.Info("devicePollerInterval", zap.Duration("interval", config.DevicePollerInterval))
	pollerJob, pollerErr := scheduler.Every(config.DevicePollerInterval).Do(devicesPoller, config, db, pollerCH)
	if pollerErr != nil {
		logger.Error("poller", zap.Any("Job", pollerJob), zap.Any("error", pollerErr))
	}
	logger.Info("deviceExportInterval", zap.Duration("interval", config.DeviceExportInterval))
	exportJob, exportErr := scheduler.Every(config.DeviceExportInterval).Do(backupScheduler, config, db, exportCH)
	if exportErr != nil {
		logger.Error("export", zap.Any("Job", exportJob), zap.Any("error", exportErr))
	}
	logger.Info("export retention job interval is 24 hours")
	exportRetentionJob, exportRetentionErr := scheduler.Every("24h").Do(rotateExports, db)
	if exportRetentionErr != nil {
		logger.Error("export", zap.Any("Job", exportRetentionJob), zap.Any("error", exportRetentionErr))
	}
	logger.Info("session cleanup job interval is 24 hours")
	sessionCleanupJob, sessionCleanupErr := scheduler.Every("24h").Do(cleanupSessions, db)
	if sessionCleanupErr != nil {
		logger.Error("session", zap.Any("Job", sessionCleanupJob), zap.Any("error", sessionCleanupErr))
	}
	scheduler.StartAsync()

	apiWorker(config, pollerCH)
	exportWorker(config, exportCH)

	wg.Wait()
}

func devicesPoller(cfg *Config, db *db.DB, pollerCH chan<- *PollerCFG) error {
	var d = &utils.Device{}

	logger.Info("starting device polling task")
	devices, err := d.GetAll(db)
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
		decryptedPw, encryptionErr := utils.DecryptString(creds.EncryptedPassword, cfg.EncryptionKey)
		if encryptionErr != nil {
			logger.Error(encryptionErr.Error())
			return err
		}
		client := &api.API{
			Address:  device.Address,
			Port:     device.ApiPort,
			Username: creds.Username,
			Password: decryptedPw,
			Async:    false,
			UseTLS:   false,
		}
		pollerCH <- &PollerCFG{Client: client, Db: db, Device: device}
	}
	return nil
}

func apiWorker(cfg *Config, pollerCH <-chan *PollerCFG) {
	logger.Info("starting MikroTik API pollers", zap.Int("count", cfg.ApiPollers))
	for x := 0; x < cfg.ApiPollers; x++ {
		go func() {
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

				if cfg.Device.Id != "" {
					dbErr = cfg.Device.Update(cfg.Db)
				} else {
					dbErr = cfg.Device.Create(cfg.Db)
				}

				if dbErr != nil {
					logger.Error(dbErr.Error())
				}
			}
		}()
	}
}

func backupScheduler(cfg *Config, db *db.DB, exportCH chan<- *BackupCFG) {
	var d = &utils.Device{}

	logger.Info("starting backup task")
	devices, err := d.GetAll(db)
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
		decryptedPw, encryptionErr := utils.DecryptString(creds.EncryptedPassword, cfg.EncryptionKey)
		if encryptionErr != nil {
			logger.Error(encryptionErr.Error())
			return
		}
		client := &ssh.SshClient{
			Host:     device.Address,
			Port:     device.SshPort,
			User:     creds.Username,
			Password: decryptedPw,
		}
		exportCH <- &BackupCFG{Client: client, Db: db, Device: device}
	}
}

func exportWorker(config *Config, exportCH <-chan *BackupCFG) {
	logger.Info("starting MikroManager export workers", zap.Int("count", config.ExportWorkers))
	for x := 0; x < config.ExportWorkers; x++ {
		wg.Add(1)
		go func() {
			for cfg := range exportCH {
				logger.Info("creating backup", zap.String("address", cfg.Client.Host))

				export, sshErr := cfg.Client.Run("/export show-sensitive")
				if sshErr == nil {
					err := s3.UploadFile(cfg.Device.Id, []byte(export))
					if err != nil {
						logger.Error(err.Error())
					}
				} else {
					logger.Error(sshErr.Error())
				}
			}
		}()
	}
}

func rotateExports(db *db.DB) {
	var err error
	var exportsList []*internal.Export
	var device *utils.Device

	logger.Info("starting exports retention task")
	err = policy.GetDefault(db)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	devices, err := device.GetAll(db)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, device := range devices {
		exports, err := s3.GetExports(device.Id)
		if err != nil {
			logger.Error(err.Error())
		} else {
			exportsList = append(exportsList, rotateHourlyExports(exports, policy.Hourly)...)
			exportsList = append(exportsList, rotateDailyExports(exports, policy.Daily)...)
			exportsList = append(exportsList, rotateWeeklyExports(exports, policy.Weekly)...)

			for _, export := range exports {
				if !exportInSlice(export, exportsList) {
					logger.Debug("deleting export", zap.String("filename", export.Key))
					err := s3.DeleteFile(export.Key)
					if err != nil {
						logger.Error(err.Error())
					}
				}
			}
		}
	}
}

func cleanupSessions(db *db.DB) {
	var err error
	var session *utils.Session

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
