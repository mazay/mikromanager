package db

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Credentials struct {
	Base
	Alias             string `gorm:"type:string;primary_key" json:"id"`
	Username          string
	EncryptedPassword string
}

func (c *Credentials) Create(db *DB) error {
	var creds Credentials
	if err := db.DB.Where("alias = ?", c.Alias).First(&creds).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return db.DB.Create(&c).Error
	}
	return fmt.Errorf("alias '%s' already exists, please pick another name", c.Alias)
}

func (c *Credentials) Update(db *DB) error {
	return db.DB.Model(&c).Where("id = ?", c.Id).Updates(&c).Error
}

func (c *Credentials) Delete(db *DB) error {
	return db.DB.Delete(&c).Error
}

func (c *Credentials) GetDefault(db *DB) error {
	return db.DB.Where("alias = ?", "Default").First(&c).Error
}

func (c *Credentials) GetById(db *DB) error {
	return db.DB.First(&c, "id = ?", c.Id).Error
}

func (c *Credentials) GetAll(db *DB) ([]*Credentials, error) {
	var creds []*Credentials
	return creds, db.DB.Find(&creds).Error
}
