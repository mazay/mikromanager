package db

import "path/filepath"

func openTestDb(dbDir string) (*DB, error) {
	dbPath := filepath.Join(dbDir, "test.db")
	db := &DB{LogLevel: "info"}
	err := db.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTestUser(db *DB) (*User, error) {
	user := &User{
		Username:          "test-user",
		EncryptedPassword: "test-password",
	}
	err := user.Create(db)
	if err != nil {
		return nil, err
	}
	return user, nil
}
