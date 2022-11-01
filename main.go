package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/mazay/mikromanager/api"
	"github.com/mazay/mikromanager/db"
	"github.com/mazay/mikromanager/http"
	"github.com/mazay/mikromanager/ssh"
	"github.com/mazay/mikromanager/utils"
	"github.com/sirupsen/logrus"
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

	policy = &utils.ExportsRetentionPolicy{Name: "Default"}
	user   = &utils.User{}

	log    = logrus.New()
	logger = log.WithFields(logrus.Fields{"app": "mikromanager"})
	wg     = sync.WaitGroup{}
	osExit = os.Exit
)

func main() {
	// Read command line args
	flag.StringVar(&configPath, "config", "config.yml", "Path to the config.yml")
	flag.StringVar(&httpPort, "http-port", "8080", "Port for the HTTP server")
	flag.Parse()

	config := readConfigFile(configPath)
	initLogger(config)

	pollerCH := make(chan *PollerCFG)
	exportCH := make(chan *BackupCFG)

	wg.Add(1)

	db := &db.DB{Path: config.DbPath}
	err = db.Init()
	if err != nil {
		logger.Panicf("DB init issue: %s", err)
		osExit(1)
	}
	defer db.Close()

	logger.Debug("ensure 'Default' exports retention policy exists")
	err = policy.GetDefault(db)
	if err != nil || policy.Id == "" {
		logger.Error(err)
	}

	if policy.Id == "" {
		policy.Hourly = 24
		policy.Daily = 14
		policy.Weekly = 26
		err = policy.Create(db)
		if err != nil {
			logger.Fatal(err)
		}
	}

	logger.Debug("ensure at least one user exists, create 'admin' otherwise")
	users, err := user.GetAll(db)
	if err != nil {
		logger.Error(err)
	}

	if len(users) == 0 {
		encryptedPw, err := utils.EncryptString("admin", config.EncryptionKey)
		if err != nil {
			logger.Error(err)
			osExit(3)
		}
		user.Username = "admin"
		user.EncryptedPassword = encryptedPw
		err = user.Create(db)
		if err != nil {
			logger.Error(err)
			osExit(3)
		}
	}

	collections, _ := db.ListCollections()
	logger.Debugf("DB has the following collections: %s", strings.Join(collections, ", "))

	go http.HttpServer("8000", db, config.EncryptionKey, config.BackupPath, logger)

	scheduler := gocron.NewScheduler(time.Local)
	logger.Infof("devicePollerInterval is %s", config.DevicePollerInterval)
	pollerJob, pollerErr := scheduler.Every(config.DevicePollerInterval).Do(devicesPoller, config, db, pollerCH)
	if pollerErr != nil {
		logger.Errorf("Job: %v, Error: %v", pollerJob, pollerErr)
	}
	logger.Infof("deviceExportInterval is %s", config.DeviceExportInterval)
	exportJob, exportErr := scheduler.Every(config.DeviceExportInterval).Do(backupScheduler, config, db, exportCH)
	if exportErr != nil {
		logger.Errorf("Job: %v, Error: %v", exportJob, exportErr)
	}
	logger.Info("export retention job interval is 90 minutes")
	exportRetentionJob, exportRetentionErr := scheduler.Every("90m").Do(rotateExports, db)
	if exportRetentionErr != nil {
		logger.Errorf("Job: %v, Error: %v", exportRetentionJob, exportRetentionErr)
	}
	logger.Info("session cleanup job interval is 24 hours")
	sessionCleanupJob, sessionCleanupErr := scheduler.Every("24h").Do(cleanupSessions, db)
	if sessionCleanupErr != nil {
		logger.Errorf("Job: %v, Error: %v", sessionCleanupJob, sessionCleanupErr)
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
		logger.Error(err)
		return err
	}
	for _, device := range devices {
		creds, err := device.GetCredentials(db)
		if err != nil {
			logger.Error(err)
			return err
		}
		logger.Debugf("using credentials '%s' for device '%s'", creds.Alias, device.Address)
		decryptedPw, encryptionErr := utils.DecryptString(creds.EncryptedPassword, cfg.EncryptionKey)
		if encryptionErr != nil {
			logger.Error(encryptionErr)
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
	logger.Infof("starting %d MikroTik API pollers", cfg.ApiPollers)
	for x := 0; x < cfg.ApiPollers; x++ {
		go func() {
			for cfg := range pollerCH {
				var fetchErr error
				var dbErr error
				logger.Infof("polling device '%s'", cfg.Client.Address)
				fetchErr = fetchResources(cfg)
				if fetchErr != nil {
					logger.Error(fetchErr)
				}

				fetchErr = fetchIdentity(cfg)
				if fetchErr != nil {
					logger.Error(fetchErr)
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
					logger.Error(dbErr)
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
		logger.Error(err)
		return
	}
	for _, device := range devices {
		creds, err := device.GetCredentials(db)
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Debugf("using credentials '%s' for device '%s'", creds.Alias, device.Address)
		decryptedPw, encryptionErr := utils.DecryptString(creds.EncryptedPassword, cfg.EncryptionKey)
		if encryptionErr != nil {
			logger.Error(encryptionErr)
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
	logger.Infof("starting %d MikroManager export workers", config.ExportWorkers)
	for x := 0; x < config.ExportWorkers; x++ {
		wg.Add(1)
		go func() {
			for cfg := range exportCH {
				logger.Infof("creating backup for device with IP address %s", cfg.Client.Host)

				export, sshErr := cfg.Client.Run("/export")
				if sshErr == nil {
					creationTime := time.Now()
					filename := fmt.Sprintf("%s/exports/%s/%d.rsc", config.BackupPath, cfg.Client.Host, creationTime.Unix())
					err := writeBackupFile(filename, export)
					if err != nil {
						logger.Error(err)
					} else {
						export := &utils.Export{
							DeviceId: cfg.Device.Id,
							Filename: filename,
						}
						err := export.Create(cfg.Db)
						if err != nil {
							logger.Fatal(err)
						}
					}
				} else {
					logger.Error(sshErr)
				}
			}
		}()
	}
}

func rotateExports(db *db.DB) {
	var err error
	var exportsList []*utils.Export
	var device *utils.Device

	logger.Info("starting exports retention task")
	err = policy.GetDefault(db)
	if err != nil {
		logger.Error(err)
		return
	}
	devices, err := device.GetAll(db)
	if err != nil {
		logger.Error(err)
		return
	}
	for _, device := range devices {
		var export *utils.Export
		exports, err := export.GetByDeviceId(db, device.Id)
		if err != nil {
			logger.Error(err)
		} else {
			exportsList = append(exportsList, rotateHourlyExports(db, exports, policy.Hourly)...)
			exportsList = append(exportsList, rotateDailyExports(db, exports, policy.Daily)...)
			exportsList = append(exportsList, rotateWeeklyExports(db, exports, policy.Weekly)...)

			for _, export := range exports {
				if !exportInSlice(export, exportsList) {
					logger.Debugf("deleting export '%s' according to the retention policy", export.Filename)
					err := export.Delete(db)
					if err != nil {
						logger.Error(err)
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
		logger.Error(err)
		return
	}

	for _, s := range sessions {
		if s.ValidThrough.Before(time.Now()) {
			logger.Debugf("session '%s' has expired, deleting", s.Id)
			err = s.Delete(db)
			if err != nil {
				logger.Error(err)
			}
		}
	}
}
