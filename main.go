package main

import (
	"log"
	"sync"
	"time"

	"github.com/mazay/mikromanager/api"
	"github.com/mazay/mikromanager/db"
	"github.com/mazay/mikromanager/http"
	"github.com/mazay/mikromanager/utils"
)

type PollerCFG struct {
	Client *api.API
	Db     *db.DB
}

var (
	wg            = sync.WaitGroup{}
	encryptionKey = "eek3eagheCo1phah4shi2Nai3ce8tiehaeVe5baph6Aixi9oorai5iepa1woh4ieQuaiz4outhakeixohn6aech8riep7beeluum"
)

func main() {
	wg.Add(1)
	pollerCH := make(chan PollerCFG)

	db := &db.DB{Path: "database.clover"}
	db.Init()
	defer db.Close()

	go http.HttpServer("8000", db, encryptionKey)

	go apiPoller(pollerCH)
	go devicesPoller(db, pollerCH)

	wg.Wait()
}

func devicesPoller(db *db.DB, pollerCH chan<- PollerCFG) {
	var d = &utils.Device{}
	log.Print("starting device poller/scheduler")
	for {
		devices, err := d.GetAll(db)
		if err != nil {
			log.Print(err)
			return
		}
		for _, device := range devices {
			creds, err := device.GetCredentials(db)
			if err != nil {
				log.Print(err)
				return
			}
			log.Printf("using credentials '%s' for device '%s'", creds.Alias, device.Address)
			decryptedPw, encryptionErr := utils.DecryptString(creds.EncryptedPassword, encryptionKey)
			if encryptionErr != nil {
				log.Print(err)
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
		time.Sleep(300 * time.Second)
	}
}

func apiPoller(pollerCH <-chan PollerCFG) {
	log.Print("starting MikroTik API poller")
	for {
		select {
		case cfg := <-pollerCH:
			log.Printf("polling device '%s'", cfg.Client.Address)
			resource, err := cfg.Client.Run("/system/resource/print")
			if err != nil {
				log.Print(err)
				return
			}
			log.Printf("fetched resource data for %s", cfg.Client.Address)

			identity, err := cfg.Client.Run("/system/identity/print")
			if err != nil {
				log.Print(err)
				return
			}
			log.Printf("identity for %s is %s", cfg.Client.Address, identity[0].Map["name"])

			values := make(map[string]interface{}, len(resource[0].Map))
			for k, v := range resource[0].Map {
				values[k] = v
			}
			values["identity"] = string(identity[0].Map["name"])
			values["polledAt"] = time.Now()
			cfg.Db.Update("devices", "address", cfg.Client.Address, values)
			// cfg.Db.Print()
			cfg.Db.Export("devices", "devices_export.json")
		}
	}
}
