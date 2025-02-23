package db

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Device struct {
	Base
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
	CredentialsId        *string           `json:"credentialsId"`
	FactorySoftware      string            `json:"factory-software"`
	FreeHddSpace         int64             `json:"free-hdd-space,string"`
	FreeMemory           int64             `json:"free-memory,string"`
	Identity             string            `json:"identity"`
	Platform             string            `json:"platform"`
	PolledAt             *time.Time        `json:"polledAt"`
	PollingSucceeded     int64             `json:"pollingSucceeded,string"`
	SshPort              string            `json:"sshPort"`
	TotalHddSpace        int64             `json:"total-hdd-space,string"`
	TotalMemory          int64             `json:"total-memory,string"`
	Uptime               string            `json:"uptime"`
	Version              string            `json:"version"`
	WriteSectSinceReboot int64             `json:"write-sect-since-reboot,string"`
	WriteSectTotal       int64             `json:"write-sect-total,string"`
	Model                string            `json:"model"`
	SerialNumber         string            `json:"serial-number"`
	FirmwareType         string            `json:"firmware-type"`
	FactoryFirmware      string            `json:"factory-firmware"`
	CurrentFirmware      string            `json:"current-firmware"`
	UpgradeFirmware      string            `json:"upgrade-firmware"`
}

func (d *Device) GetAll(db *DB) ([]*Device, error) {
	var deviceList []*Device
	return deviceList, db.DB.Find(&deviceList).Error
}

func (d *Device) GetCredentials(db *DB) (*Credentials, error) {
	var c = &Credentials{}
	if *d.CredentialsId == "" {
		return c, c.GetDefault(db)
	} else {
		c.ID = *d.CredentialsId
		return c, c.GetById(db)
	}
}

func (d *Device) Create(db *DB) error {
	var dev Device

	d.PollingSucceeded = -1
	// check if device with given address already exist
	if err := db.DB.Where("address = ?", d.Address).Find(&dev).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return db.DB.Create(d).Error
	}

	return fmt.Errorf("Device with address '%s' already exists", d.Address)
}

func (d *Device) Update(db *DB) error {
	var dev Device

	if err := db.DB.Where("address = ?", d.Address).Find(&dev).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return d.Create(db)
	}

	return db.DB.Model(d).Updates(d).Error
}

func (d *Device) GetById(db *DB) error {
	return db.DB.First(d, "ID = ?", d.ID).Error
}

func (d *Device) Delete(db *DB) error {
	return db.DB.Delete(d).Error
}
