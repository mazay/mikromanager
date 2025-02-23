package db

import (
	"time"
)

type Device struct {
	Base
	Address              string `gorm:"unique"`
	ApiPort              string
	ArchitectureName     string
	BadBlocks            int64
	BoardName            string
	BuildTime            string
	CPU                  string
	CpuCount             int64
	CpuFrequency         int64
	CpuLoad              int64
	CredentialsId        string
	FactorySoftware      string
	FreeHddSpace         int64
	FreeMemory           int64
	Identity             string
	Platform             string
	PolledAt             time.Time
	PollingSucceeded     int64
	SshPort              string
	TotalHddSpace        int64
	TotalMemory          int64
	Uptime               string
	Version              string
	WriteSectSinceReboot int64
	WriteSectTotal       int64
	Model                string
	SerialNumber         string
	FirmwareType         string
	FactoryFirmware      string
	CurrentFirmware      string
	UpgradeFirmware      string
}

func (d *Device) GetAll(db *DB) ([]*Device, error) {
	var deviceList []*Device
	return deviceList, db.DB.Find(&deviceList).Error
}

func (d *Device) GetCredentials(db *DB) (*Credentials, error) {
	var c = &Credentials{}
	if d.CredentialsId == "" {
		return c, c.GetDefault(db)
	} else {
		c.Id = d.CredentialsId
		return c, c.GetById(db)
	}
}

func (d *Device) Create(db *DB) error {
	d.PollingSucceeded = -1
	return db.DB.Create(&d).Error
}

func (d *Device) Update(db *DB) error {
	return db.DB.Model(&d).Where("id = ?", d.Id).Updates(d).Error
}

func (d *Device) GetById(db *DB) error {
	return db.DB.First(&d, "id = ?", d.Id).Error
}

func (d *Device) Delete(db *DB) error {
	return db.DB.Delete(&d).Error
}
