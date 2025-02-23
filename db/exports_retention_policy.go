package db

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ExportsRetentionPolicy struct {
	Base
	Name   string `json:"name"`
	Hourly int64  `json:"hourly"`
	Daily  int64  `json:"daily"`
	Weekly int64  `json:"weekly"`
}

func (rp *ExportsRetentionPolicy) Create(db *DB) error {
	var p ExportsRetentionPolicy

	// check if policy with that name already exist
	if err := db.DB.Where("name = ?", rp.Name).Find(&p).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return db.DB.Create(rp).Error
	}

	return fmt.Errorf("retention poilicy '%s' already exists, please pick another name", rp.Name)
}

func (rp *ExportsRetentionPolicy) Update(db *DB) error {
	var p ExportsRetentionPolicy
	if err := db.DB.Where("ID = ?", rp.ID).Find(&p).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return rp.Create(db)
	}

	return db.DB.Model(rp).Updates(rp).Error
}

func (rp *ExportsRetentionPolicy) GetDefault(db *DB) error {
	if err := db.DB.Where("name = ?", "Default").Find(rp).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		rp.Hourly = 24
		rp.Daily = 14
		rp.Weekly = 26
		return rp.Create(db)
	}

	return nil
}
