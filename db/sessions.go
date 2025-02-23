package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Session struct {
	Base
	UserId       string    `json:"userId"`
	ValidThrough time.Time `json:"expires"`
}

func (s *Session) Expired() bool {
	if s.ValidThrough.IsZero() {
		return true
	}
	return s.ValidThrough.Before(time.Now())
}

func (s *Session) Create(db *DB) error {
	return db.DB.Create(s).Error
}

func (s *Session) Delete(db *DB) error {
	return db.DB.Delete(s).Error
}

func (s *Session) GetById(db *DB) error {
	return db.DB.First(s, "ID = ?", s.ID).Error
}

func (s *Session) Update(db *DB) error {
	var session Session

	if err := db.DB.Where("ID = ?", s.ID).Find(&session).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return s.Create(db)
	}

	return db.DB.Model(s).Updates(s).Error
}

func (s *Session) GetAll(db *DB) ([]*Session, error) {
	var sessionList []*Session
	return sessionList, db.DB.Find(&sessionList).Error
}

func (s *Session) GetByUserId(db *DB, userId string) ([]*Session, error) {
	var sessionList []*Session
	return sessionList, db.DB.Where("userId = ?", s.UserId).Find(&sessionList).Error
}
