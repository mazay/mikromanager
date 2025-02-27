package db

import (
	"errors"

	"gorm.io/gorm"
)

type Health struct {
	Id                string  `gorm:"type:string;primary_key"`
	Voltage           float32 `json:"voltage"`
	CpuTemp           float32 `json:"cpu-temperature"`
	BoardTemp1        float32 `json:"board-temperature1"`
	BoardTemp2        float32 `json:"board-temperature2"`
	SfpTemp           float32 `json:"sfp-temperature"`
	FanState          string  `json:"fan-state"`
	FanSpeed          int     `json:"fan-speed"`
	Psu1Voltage       float32 `json:"psu1-voltage"`
	Psu2Voltage       float32 `json:"psu2-voltage"`
	PoeOutConsumption float32 `json:"poe-out-consumption"`
	JackVoltage       float32 `json:"jack-voltage"`
	TwoPinVoltage     float32 `json:"2pin-voltage"`
	PoeInVoltage      float32 `json:"poe-in-voltage"`
	DeviceId          string
	Device            *Device
}

// BeforeCreate will set the ID of the Health object to the ID of its associated Device before creation.
func (h *Health) BeforeCreate(tx *gorm.DB) error {
	h.Id = h.DeviceId
	return nil
}

// Save will create a new health entry in the database if it does not already exist for the given device,
// or update the existing entry if it does. It returns an error if the save operation fails.
func (h *Health) Save(db *DB) error {
	if err := db.DB.Where("id = ?", h.DeviceId).First(&h).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return db.DB.Create(&h).Error
	}
	return db.DB.Save(&h).Where("id = ?", h.DeviceId).Error
}

// Delete will delete the health entry from the database that matches the current
// object's DeviceId. It returns an error if the deletion fails.
func (h *Health) Delete(db *DB) error {
	return db.DB.Delete(&h).Error
}

// GetByDeviceId retrieves a health entry from the database using the current object's DeviceId and
// populates the current object with its values. It returns nil if the retrieval is successful,
// or an error if the retrieval fails. If the health entry does not exist, the function returns nil.
func (h *Health) GetByDeviceId(db *DB) error {
	err := db.DB.Where("device_id = ?", h.DeviceId).First(&h).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}
