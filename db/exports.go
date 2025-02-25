package db

import (
	"time"

	"gorm.io/gorm/clause"
)

type Export struct {
	Base
	S3Key        string `json:"s3Key"`
	LastModified *time.Time
	ETag         string
	Size         *int64
	DeviceId     string
	Device       *Device
}

func (e *Export) Save(db *DB) error {
	return db.DB.Save(&e).Error
}

func (e *Export) Delete(db *DB) error {
	return db.DB.Delete(&e).Error
}

func (e *Export) GetAll(db *DB) ([]*Export, error) {
	var exportList []*Export
	return exportList, db.DB.Preload(clause.Associations).Find(&exportList).Error
}

func (e *Export) GetById(db *DB) error {
	return db.DB.Preload(clause.Associations).First(&e, "id = ?", e.Id).Error
}

func (e *Export) GetByDeviceId(db *DB, deviceId string) ([]*Export, error) {
	var exportList []*Export
	return exportList, db.DB.Preload(clause.Associations).Find(&exportList, "device_id = ?", deviceId).Error
}
