package utils

import (
	"encoding/json"
	"fmt"
	"time"

	database "github.com/mazay/mikromanager/db"
)

type Device struct {
	Id                   string            `json:"_id"`
	Address              string            `json:"address"`
	ApiPort              string            `json:"apiPort"`
	ArchitectureName     string            `json:"architecture-name"`
	BadBlocks            int64             `json:"bad-blocks,string"`
	BoardName            DeviceBoardName   `json:"board-name"`
	BuildTime            FirmwareBuildTime `json:"build-time"`
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
	SshPort              string            `json:"sshPort"`
	TotalHddSpace        int64             `json:"total-hdd-space,string"`
	TotalMemory          int64             `json:"total-memory,string"`
	Updated              time.Time         `json:"updated"`
	Uptime               string            `json:"uptime"`
	Version              string            `json:"version"`
	WriteSectSinceReboot int64             `json:"write-sect-since-reboot,string"`
	WriteSectTotal       int64             `json:"write-sect-total,string"`
	Model                string            `json:"model"`
	SerialNumber         string            `json:"serial-number"`
	FirmwareType         string            `json:"firmware-type"`
}

func (d *Device) GetAll(db *database.DB) ([]*Device, error) {
	var deviceList []*Device

	db.Sort("address", 1)
	docs, err := db.FindAll(db.Collections["devices"])
	if err != nil {
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
	// check if device with given address already exist
	exists, _ := db.Exists(db.Collections["devices"], "address", d.Address)
	if exists {
		return fmt.Errorf("Device with address '%s' already exists", d.Address)
	}
	d.PollingSucceeded = -1
	d.Created = time.Now()
	d.Updated = time.Now()
	inrec, err := json.Marshal(d)
	if err != nil {
		return err
	}
	err = json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return err
	}
	d.Id, err = db.Insert(db.Collections["devices"], inInterface)
	return err
}

func (d *Device) Update(db *database.DB) error {
	var inInterface map[string]interface{}
	d.Updated = time.Now()
	inrec, err := json.Marshal(d)
	if err != nil {
		return err
	}
	err = json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return err
	}
	return db.UpdateById(db.Collections["devices"], d.Id, inInterface)
}

func (d *Device) GetById(db *database.DB) error {
	doc, err := db.FindById(db.Collections["devices"], d.Id)
	if err != nil {
		return err
	}

	dj, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dj, d)
	if err != nil {
		return err
	}

	return err
}

func (d *Device) Delete(db *database.DB) error {
	var err error

	// delete exports first
	e := &Export{}
	exports, err := e.GetByDeviceId(db, d.Id)
	if err != nil {
		return err
	}

	for _, export := range exports {
		err = export.Delete(db)
		if err != nil {
			return err
		}
	}

	// finally delete the device
	return db.DeleteById(db.Collections["devices"], d.Id)
}
