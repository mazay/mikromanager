package db

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Credentials struct {
	Base
	Alias             string `json:"alias"`
	Username          string `json:"username"`
	EncryptedPassword string `json:"encryptedPassword"`
}

func (c *Credentials) Create(db *DB) error {
	var creds Credentials
	if err := db.DB.Where("alias = ?", c.Alias).Find(&creds).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return db.DB.Create(c).Error
	}
	return fmt.Errorf("alias '%s' already exists, please pick another name", c.Alias)
}

func (c *Credentials) Update(db *DB) error {
	var creds Credentials
	if err := db.DB.Where("alias = ?", c.Alias).Find(&creds).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Create(db)
	}
	return db.DB.Model(c).Updates(c).Error
}

func (c *Credentials) Delete(db *DB) error {
	return db.DB.Delete(c).Error
}

func (c *Credentials) GetDefault(db *DB) error {
	return db.DB.Where("alias = ?", "Default").First(c).Error
}

func (c *Credentials) GetById(db *DB) error {
	return db.DB.First(c, "ID = ?", c.ID).Error
}

func (c *Credentials) GetAll(db *DB) ([]*Credentials, error) {
	var creds []*Credentials
	return creds, db.DB.Find(&creds).Error
}
