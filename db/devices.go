package db

import (
	"time"
)

type Device struct {
	Base
	Address              string `gorm:"unique" json:"address"`
	ApiPort              string `json:"apiPort"`
	ArchitectureName     string `json:"architecture-name"`
	BadBlocks            int64  `json:"bad-blocks,string"`
	BoardName            string `json:"board-name"`
	BuildTime            string `json:"build-time"`
	CPU                  string `json:"cpu"`
	CpuCount             int64  `json:"cpu-count,string"`
	CpuFrequency         int64  `json:"cpu-frequency,string"`
	CpuLoad              int64  `json:"cpu-load,string"`
	CredentialsId        string
	FactorySoftware      string `json:"factory-software"`
	FreeHddSpace         int64  `json:"free-hdd-space,string"`
	FreeMemory           int64  `json:"free-memory,string"`
	Identity             string `json:"identity"`
	Platform             string `json:"platform"`
	PolledAt             time.Time
	PollingSucceeded     int64
	SshPort              string
	TotalHddSpace        int64  `json:"total-hdd-space,string"`
	TotalMemory          int64  `json:"total-memory,string"`
	Uptime               string `json:"uptime"`
	Version              string `json:"version"`
	WriteSectSinceReboot int64  `json:"write-sect-since-reboot,string"`
	WriteSectTotal       int64  `json:"write-sect-total,string"`
	Model                string `json:"model"`
	SerialNumber         string `json:"serial-number"`
	FirmwareType         string `json:"firmware-type"`
	FactoryFirmware      string `json:"factory-firmware"`
	CurrentFirmware      string `json:"current-firmware"`
	UpgradeFirmware      string `json:"upgrade-firmware"`
}

// GetAll fetches all devices from the database.
//
// The function returns a slice of *Device instances and an error.
func (d *Device) GetAll(db *DB) ([]*Device, error) {
	var deviceList []*Device
	return deviceList, db.DB.Find(&deviceList).Error
}

// GetCredentials returns the credentials object associated with the device,
// or the default credential set if the device has not been configured to use
// a specific set of credentials.
func (d *Device) GetCredentials(db *DB) (*Credentials, error) {
	var c = &Credentials{}
	if d.CredentialsId == "" {
		return c, c.GetDefault(db)
	} else {
		c.Id = d.CredentialsId
		return c, c.GetById(db)
	}
}

// Create will create a new device entry in the database with the current object's values.
// The function automatically sets the PollingSucceeded field to -1 to indicate that the
// device has not been polled yet.
//
// The function returns an error if the creation fails.
func (d *Device) Create(db *DB) error {
	d.PollingSucceeded = -1
	return db.DB.Create(&d).Error
}

// Update will update an existing device entry in the database with the current
// object's values. It returns an error if the update fails.
func (d *Device) Update(db *DB) error {
	return db.DB.Model(&d).Where("id = ?", d.Id).Updates(d).Error
}

// GetById fetches a device entry from the database using the current object's ID
// and populates the current object with its values. It returns an error if the fetch
// fails.
func (d *Device) GetById(db *DB) error {
	return db.DB.First(&d, "id = ?", d.Id).Error
}

// Delete will delete an existing device entry from the database that matches the
// current object's ID. It returns an error if the deletion fails.
func (d *Device) Delete(db *DB) error {
	return db.DB.Delete(&d).Error
}
