package db

import (
	"errors"

	"gorm.io/gorm"
)

type ExportsRetentionPolicy struct {
	Base
	Name   string `gorm:"unique"`
	Hourly int64
	Daily  int64
	Weekly int64
}

func (rp *ExportsRetentionPolicy) Create(db *DB) error {
	return db.DB.Create(&rp).Error
}

func (rp *ExportsRetentionPolicy) Update(db *DB) error {
	return db.DB.Model(&rp).Where("id = ?", rp.Id).Updates(rp).Error
}

func (rp *ExportsRetentionPolicy) GetDefault(db *DB) error {
	if err := db.DB.Where("name = ?", "Default").First(&rp).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		rp.Hourly = 24
		rp.Daily = 14
		rp.Weekly = 26
		return rp.Create(db)
	}

	return nil
}
