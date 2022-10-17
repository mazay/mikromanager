package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mazay/mikromanager/api"
	"github.com/mazay/mikromanager/db"
	"github.com/mazay/mikromanager/http"
	"github.com/mazay/mikromanager/utils"
	"github.com/sirupsen/logrus"
)

type PollerCFG struct {
	Client *api.API
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

	go http.HttpServer("8000", db, config.EncryptionKey, logger)

	apiPoller(config, pollerCH)
	go devicesPoller(config, db, pollerCH)

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
				Port:     device.Port,
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
	logger.Infof("starting %s MikroTik API pollers", strconv.Itoa(cfg.ApiPollers))
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

func fetchResources(cfg *PollerCFG) error {
	resource, err := cfg.Client.Run("/system/resource/print")
	if err != nil {
		return err
	}
	logger.Debugf("fetched resource data for %s", cfg.Client.Address)
	inrec, _ := json.Marshal(resource[0].Map)
	return json.Unmarshal(inrec, &cfg.Device)
}

func fetchIdentity(cfg *PollerCFG) error {
	identity, err := cfg.Client.Run("/system/identity/print")
	if err != nil {
		return err
	}
	if len(identity) > 0 {
		logger.Debugf("identity for %s is %s", cfg.Client.Address, identity[0].Map["name"])
		cfg.Device.Identity = string(identity[0].Map["name"])
		return nil
	}
	return fmt.Errorf("got an empty identity data")
}
