package db

import (
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
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
	if base.Id == "" {
		base.Id = uuid.New().String()
	}
	return nil
}

type DB struct {
	DB       *gorm.DB
	LogLevel string
}

// Open initializes the database connection using the provided file path
// and configures the logging level. It also migrates the schema for the
// specified models. Returns an error if the connection or migration fails.
func (db *DB) Open(path string) error {
	logLevel, ok := dbLogLevels[db.LogLevel]
	if !ok {
		return fmt.Errorf("invalid log level: %s", db.LogLevel)
	}
	if path == "" {
		return fmt.Errorf("path to the database is empty")
	}
	gormDB, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return err
	}

	db.DB = gormDB

	// Migrate the schemas
	err = db.DB.AutoMigrate(
		&Credentials{},
		&User{},
		&Device{},
		&ExportsRetentionPolicy{},
		&Session{},
		&DeviceGroup{},
		&Export{},
		&Health{},
	)
	if err != nil {
		return err
	}

	return nil
}

// Close the underlying database connection.
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
