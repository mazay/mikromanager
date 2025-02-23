package db

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	Base
	Username          string `json:"username"`
	EncryptedPassword string `json:"encryptedPassword"`
}

func (u *User) Create(db *DB) error {
	var creds Credentials
	if err := db.DB.Where("username = ?", u.Username).Find(&creds).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return db.DB.Debug().Create(u).Error
	}
	return nil
}

func (u *User) Delete(db *DB) error {
	return db.DB.Delete(u).Error
}

func (u *User) GetById(db *DB) error {
	return db.DB.First(u, "ID = ?", u.ID).Error
}

func (u *User) GetByUsername(db *DB) error {
	return db.DB.First(u, "username = ?", u.Username).Error
}

func (u *User) Update(db *DB) error {
	var user User
	if err := db.DB.Where("ID = ?", u.ID).Find(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return u.Create(db)
	}
	return db.DB.Model(u).Updates(u).Error
}

func (u *User) GetAll(db *DB) ([]*User, error) {
	var userList []*User
	return userList, db.DB.Find(&userList).Error
}
