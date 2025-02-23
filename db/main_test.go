package db

import (
	"path/filepath"
	"testing"
)

func TestDBOpen(t *testing.T) {
	// Test successful database connection and migration
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db := &DB{LogLevel: "silent"}
	err := db.Open(dbPath)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if db.DB == nil {
		t.Errorf("expected db.DB to be set")
	}
}

func TestDBOpenInvalidPath(t *testing.T) {
	// Test error handling for invalid database path
	db := &DB{LogLevel: "info"}
	err := db.Open("")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestDBOpenInvalidLogLevel(t *testing.T) {
	// Test error handling for invalid log level
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db := &DB{LogLevel: "invalid"}
	err := db.Open(dbPath)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
