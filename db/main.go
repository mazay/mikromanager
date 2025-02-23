package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbLogLevels = map[string]logger.LogLevel{
		"silent": logger.Silent,
		"error":  logger.Error,
		"warn":   logger.Warn,
		"info":   logger.Info,
	}
)

// Base contains common columns for all tables.
type Base struct {
	Id        string `gorm:"type:string;primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	base.Id = uuid.New().String()
	return nil
}

type DB struct {
	DB       *gorm.DB
	LogLevel string
}

func (db *DB) Open(path string) error {
	gormDB, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(dbLogLevels[db.LogLevel]),
	})
	if err != nil {
		return err
	}

	db.DB = gormDB

	// Migrate the schema
	err = db.DB.AutoMigrate(&Credentials{}, &User{}, &Device{}, &ExportsRetentionPolicy{}, &Session{})
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
