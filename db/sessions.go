package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	Base
	UserId       string
	ValidThrough time.Time
}

// BeforeCreate will set the session to expire in 24 hours
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	s.Id = uuid.New().String()
	s.ValidThrough = time.Now().Add(time.Hour * 24)
	return nil
}

func (s *Session) Expired() bool {
	if s.ValidThrough.IsZero() {
		return true
	}
	return s.ValidThrough.Before(time.Now())
}

func (s *Session) Create(db *DB) error {
	return db.DB.Create(&s).Error
}

func (s *Session) Delete(db *DB) error {
	return db.DB.Delete(&s).Error
}

func (s *Session) GetById(db *DB) error {
	return db.DB.First(&s, "id = ?", s.Id).Error
}

func (s *Session) Update(db *DB) error {
	return db.DB.Model(&s).Where("id = ?", s.Id).Updates(s).Error
}

func (s *Session) GetAll(db *DB) ([]*Session, error) {
	var sessionList []*Session
	return sessionList, db.DB.Find(&sessionList).Error
}

func (s *Session) GetByUserId(db *DB, userId string) ([]*Session, error) {
	var sessionList []*Session
	return sessionList, db.DB.Where("userId = ?", s.UserId).Find(&sessionList).Error
}
