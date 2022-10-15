package main

import (
	"flag"
	"os"
	"strconv"
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

	pollerCH := make(chan PollerCFG)

	wg.Add(1)

	db := &db.DB{Path: config.DbPath}
	db.Init()
	defer db.Close()

	go http.HttpServer("8000", db, config.EncryptionKey, logger)

	apiPoller(config, pollerCH)
	go devicesPoller(config, db, pollerCH)

	wg.Wait()
}

func devicesPoller(cfg *Config, db *db.DB, pollerCH chan<- PollerCFG) {
	var d = &utils.Device{}
	logger.Info("starting device poller/scheduler")
	logger.Infof("apiPollerInterval is %s", cfg.ApiPollerInterval)
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
				logger.Error(err)
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
			pollerCH <- PollerCFG{Client: client, Db: db}
		}
		time.Sleep(cfg.ApiPollerInterval)
	}
}

func apiPoller(cfg *Config, pollerCH <-chan PollerCFG) {
	logger.Infof("starting %s MikroTik API pollers", strconv.Itoa(cfg.ApiPollers))
	for x := 0; x < cfg.ApiPollers; x++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case cfg := <-pollerCH:
					logger.Infof("polling device '%s'", cfg.Client.Address)
					resource, err := cfg.Client.Run("/system/resource/print")
					if err != nil {
						logger.Error(err)
						recover()
					} else {
						logger.Debugf("fetched resource data for %s", cfg.Client.Address)
					}

					identity, err := cfg.Client.Run("/system/identity/print")
					if err != nil {
						logger.Error(err)
						recover()
					} else {
						logger.Debugf("identity for %s is %s", cfg.Client.Address, identity[0].Map["name"])

						values := make(map[string]interface{}, len(resource[0].Map))
						for k, v := range resource[0].Map {
							values[k] = v
						}
						values["identity"] = string(identity[0].Map["name"])
						values["polledAt"] = time.Now()
						cfg.Db.Update("devices", "address", cfg.Client.Address, values)
					}
				}
			}
		}()
	}
}
