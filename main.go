package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

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
	configPath string
	httpPort   string

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
	err := db.Init()
	if err != nil {
		logger.Panicf("DB init issue: %s", err)
		osExit(1)
	}
	defer db.Close()

	collections, _ := db.ListCollections()
	logger.Debugf("DB has the following collections: %s", strings.Join(collections, ", "))

	go http.HttpServer("8000", db, config.EncryptionKey, config.BackupPath, logger)

	apiPoller(config, pollerCH)
	exportWorker(config, exportCH)
	go devicesPoller(config, db, pollerCH)
	go backupScheduler(config, db, exportCH)

	wg.Wait()
}

func devicesPoller(cfg *Config, db *db.DB, pollerCH chan<- *PollerCFG) {
	var d = &utils.Device{}
	wg.Add(1)
	logger.Info("starting device poller/scheduler")
	logger.Debugf("devicePollerInterval is %s", cfg.DevicePollerInterval)
	for {
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
		time.Sleep(cfg.DevicePollerInterval)
	}
}

func apiPoller(cfg *Config, pollerCH <-chan *PollerCFG) {
	logger.Infof("starting %d MikroTik API pollers", cfg.ApiPollers)
	for x := 0; x < cfg.ApiPollers; x++ {
		wg.Add(1)
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
	wg.Add(1)
	logger.Info("starting backup scheduler")
	logger.Debugf("deviceExportInterval is %s", cfg.DeviceExportInterval)
	for {
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
		time.Sleep(cfg.DeviceExportInterval)
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
