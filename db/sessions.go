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

// BeforeCreate will set the session to expire in 24 hours after creation and set a UUID rather than numeric ID.
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	s.Id = uuid.New().String()
	if s.ValidThrough.IsZero() {
		s.ValidThrough = time.Now().Add(time.Hour * 24)
	}
	return nil
}

// Expired will check if the session is expired or not. If the session ValidThrough value is zero, it is considered expired.
func (s *Session) Expired() bool {
	if s.ValidThrough.IsZero() {
		return true
	}
	return s.ValidThrough.Before(time.Now())
}

// Create will create a new session entry in the database with the current
// object's values. It returns an error if the creation fails.
func (s *Session) Create(db *DB) error {
	return db.DB.Create(&s).Error
}

// Delete will delete an existing session entry from the database that
// matches the current object's ID. It returns an error if the deletion fails.
func (s *Session) Delete(db *DB) error {
	return db.DB.Delete(&s).Error
}

// GetById fetches a session entry from the database using the current object's ID
// and populates the current object with its values. It returns an error if the fetch
// fails.
func (s *Session) GetById(db *DB) error {
	return db.DB.First(&s, "id = ?", s.Id).Error
}

// Update will update an existing session entry in the database with the current
// object's values. It returns an error if the update fails.
func (s *Session) Update(db *DB) error {
	return db.DB.Model(&s).Where("id = ?", s.Id).Updates(s).Error
}

// GetAll retrieves all sessions entries from the database and returns them
// as a slice of Session pointers. It returns an error if the retrieval fails.
func (s *Session) GetAll(db *DB) ([]*Session, error) {
	var sessionList []*Session
	return sessionList, db.DB.Find(&sessionList).Error
}

// GetByUserId retrieves all session entries from the database associated with a given user ID.
// It returns a slice of Session pointers and an error if the retrieval fails.
func (s *Session) GetByUserId(db *DB, userId string) ([]*Session, error) {
	var sessionList []*Session
	return sessionList, db.DB.Where("user_id = ?", s.UserId).Find(&sessionList).Error
}
