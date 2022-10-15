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
	Id                   string    `json:"_id"`
	Address              string    `json:"address"`
	ArchitectureName     string    `json:"architecture-name"`
	BadBlocks            int64     `json:"bad-blocks"`
	BoardName            string    `json:"board-name"`
	BuildTime            string    `json:"build-time"`
	CPU                  string    `json:"cpu"`
	CpuCount             int       `json:"cpu-count"`
	CpuFrequency         int       `json:"cpu-frequency"`
	CpuLoad              int       `json:"cpu-load"`
	Created              time.Time `json:"created"`
	CredentialsId        string    `json:"credentialsId"`
	FactorySoftware      string    `json:"factory-software"`
	FreeHddSpace         int64     `json:"free-hdd-space"`
	FreeMemory           int64     `json:"free-memory"`
	Identity             string    `json:"identity"`
	PolledAt             time.Time `json:"polledAt"`
	Platform             string    `json:"platform"`
	TotalHddSpace        int64     `json:"total-hdd-space"`
	TotalMemory          int64     `json:"total-memory"`
	Updated              time.Time `json:"updated"`
	Uptime               string    `json:"uptime"`
	Version              string    `json:"version"`
	WriteSectSinceReboot int64     `json:"write-sect-since-reboot"`
	WriteSectTotal       int64     `json:"write-sect-total"`
}

func (d *Device) GetAll(db *database.DB) ([]*Device, error) {
	var deviceList []*Device

	docs, err := db.FindAll("devices")
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
	exists, _ := db.Exists("devices", "address", d.Address)
	if exists {
		return errors.New(fmt.Sprintf("Device with address '%s' already exists", d.Address))
	}
	d.Created = time.Now()
	d.Updated = time.Now()
	inrec, _ := json.Marshal(d)
	json.Unmarshal(inrec, &inInterface)
	_, err := db.Insert("devices", inInterface)
	return err
}

func (d *Device) Update(db *database.DB) error {
	var inInterface map[string]interface{}
	d.Updated = time.Now()
	inrec, _ := json.Marshal(d)
	json.Unmarshal(inrec, &inInterface)
	return db.Update("devices", "_id", d.Id, inInterface)
}

func (d *Device) GetById(db *database.DB) error {
	doc, err := db.FindById("devices", d.Id)
	if err != nil {
		log.Fatal(err)
		return err
	}

	dj, err := json.Marshal(doc)
	json.Unmarshal(dj, d)

	return err
}

func (d *Device) Delete(db *database.DB) error {
	return db.DeleteById("devices", d.Id)
}
