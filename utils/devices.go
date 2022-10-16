package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	database "github.com/mazay/mikromanager/db"
)

type Device struct {
	Id                   string            `json:"_id"`
	Address              string            `json:"address"`
	ArchitectureName     string            `json:"architecture-name"`
	BadBlocks            int64             `json:"bad-blocks,string"`
	BoardName            string            `json:"board-name"`
	BuildTime            FirmwareBuildTime `json:"build-time,string"`
	CPU                  string            `json:"cpu"`
	CpuCount             int64             `json:"cpu-count,string"`
	CpuFrequency         int64             `json:"cpu-frequency,string"`
	CpuLoad              int64             `json:"cpu-load,string"`
	Created              time.Time         `json:"created"`
	CredentialsId        string            `json:"credentialsId"`
	FactorySoftware      string            `json:"factory-software"`
	FreeHddSpace         int64             `json:"free-hdd-space,string"`
	FreeMemory           int64             `json:"free-memory,string"`
	Identity             string            `json:"identity"`
	Platform             string            `json:"platform"`
	PolledAt             time.Time         `json:"polledAt"`
	PollingSucceeded     int64             `json:"pollingSucceeded,string"`
	Port                 string            `json:"port"`
	TotalHddSpace        int64             `json:"total-hdd-space,string"`
	TotalMemory          int64             `json:"total-memory,string"`
	Updated              time.Time         `json:"updated"`
	Uptime               string            `json:"uptime"`
	Version              string            `json:"version"`
	WriteSectSinceReboot int64             `json:"write-sect-since-reboot,string"`
	WriteSectTotal       int64             `json:"write-sect-total,string"`
}

func (d *Device) GetAll(db *database.DB) ([]*Device, error) {
	var deviceList []*Device

	docs, err := db.FindAll(db.Collections["devices"])
	if err != nil {
		log.Fatal(err)
		return deviceList, err
	}

	for _, doc := range docs {
		dm := &Device{}
		dj, _ := json.Marshal(doc)
		_ = json.Unmarshal(dj, dm)
		deviceList = append(deviceList, dm)
	}

	return deviceList, nil
}

func (d *Device) GetCredentials(db *database.DB) (*Credentials, error) {
	var c = &Credentials{}
	if d.CredentialsId == "" {
		err := c.GetDefault(db)
		return c, err
	} else {
		c.Id = d.CredentialsId
		err := c.GetById(db)
		return c, err
	}
}

func (d *Device) Create(db *database.DB) error {
	var inInterface map[string]interface{}
	// check if credentials with that alias already exist
	exists, _ := db.Exists(db.Collections["devices"], "address", d.Address)
	if exists {
		return errors.New(fmt.Sprintf("Device with address '%s' already exists", d.Address))
	}
	d.Created = time.Now()
	d.Updated = time.Now()
	inrec, _ := json.Marshal(d)
	json.Unmarshal(inrec, &inInterface)
	_, err := db.Insert(db.Collections["devices"], inInterface)
	return err
}

func (d *Device) Update(db *database.DB) error {
	var inInterface map[string]interface{}
	d.Updated = time.Now()
	inrec, _ := json.Marshal(d)
	json.Unmarshal(inrec, &inInterface)
	return db.UpdateById(db.Collections["devices"], d.Id, inInterface)
}

func (d *Device) GetById(db *database.DB) error {
	doc, err := db.FindById(db.Collections["devices"], d.Id)
	if err != nil {
		log.Fatal(err)
		return err
	}

	dj, err := json.Marshal(doc)
	json.Unmarshal(dj, d)

	return err
}

func (d *Device) Delete(db *database.DB) error {
	return db.DeleteById(db.Collections["devices"], d.Id)
}
