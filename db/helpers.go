package db

import "path/filepath"

// openTestDb creates a test database at the given directory and returns a pointer to a DB connection
// object. This function is used for unit testing purposes.
func openTestDb(dbDir string) (*DB, error) {
	dbPath := filepath.Join(dbDir, "test.db")
	db := &DB{LogLevel: "info"}
	err := db.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// createTestUser creates and inserts a test user into the provided database.
// It returns a pointer to the created User object and an error, if any occurs
// during the user creation process. This function is used for unit testing purposes.
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
