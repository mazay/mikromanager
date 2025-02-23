package db

type User struct {
	Base
	Username          string `gorm:"unique"`
	EncryptedPassword string
}

// Create will create a new user entry in the database with the current
// object's values. It returns an error if the creation fails.
func (u *User) Create(db *DB) error {
	return db.DB.Create(&u).Error
}

// Delete will delete an existing user entry from the database that
// matches the current object's ID. It returns an error if the deletion fails.
func (u *User) Delete(db *DB) error {
	return db.DB.Delete(&u).Error
}

// GetById fetches a user entry from the database using the current object's ID
// and populates the current object with its values. It returns an error if the
// fetch fails.
func (u *User) GetById(db *DB) error {
	return db.DB.First(&u, "id = ?", u.Id).Error
}

// GetByUsername fetches a user entry from the database using the current object's
// username and populates the current object with its values. It returns an error if
// the fetch fails.
func (u *User) GetByUsername(db *DB) error {
	return db.DB.First(&u, "username = ?", u.Username).Error
}

// Update will update an existing user entry in the database with the current
// object's values. It returns an error if the update fails.
func (u *User) Update(db *DB) error {
	return db.DB.Model(&u).Where("id = ?", u.Id).Updates(u).Error
}

// GetAll retrieves all user entries from the database and returns them
// as a slice of *User instances. It returns an error if the retrieval fails.
func (u *User) GetAll(db *DB) ([]*User, error) {
	var userList []*User
	return userList, db.DB.Find(&userList).Error
}
