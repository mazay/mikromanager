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
	BadBlocks            string    `json:"bad-blocks"`
	BoardName            string    `json:"board-name"`
	BuildTime            string    `json:"build-time"`
	CPU                  string    `json:"cpu"`
	CpuCount             string    `json:"cpu-count"`
	CpuFrequency         string    `json:"cpu-frequency"`
	CpuLoad              string    `json:"cpu-load"`
	Created              time.Time `json:"created"`
	CredentialsId        string    `json:"credentialsId"`
	FactorySoftware      string    `json:"factory-software"`
	FreeHddSpace         string    `json:"free-hdd-space"`
	FreeMemory           string    `json:"free-memory"`
	Identity             string    `json:"identity"`
	Platform             string    `json:"platform"`
	PolledAt             time.Time `json:"polledAt"`
	PollingSucceeded     string    `json:"pollingSucceeded"`
	Port                 string    `json:"port"`
	TotalHddSpace        string    `json:"total-hdd-space"`
	TotalMemory          string    `json:"total-memory"`
	Updated              time.Time `json:"updated"`
	Uptime               string    `json:"uptime"`
	Version              string    `json:"version"`
	WriteSectSinceReboot string    `json:"write-sect-since-reboot"`
	WriteSectTotal       string    `json:"write-sect-total"`
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
