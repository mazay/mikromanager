package db

type Credentials struct {
	Base
	Alias             string `gorm:"unique"`
	Username          string
	EncryptedPassword string
}

// Create will create a new credentials entry in the database with the current
// object's values. It returns an error if the creation fails.
func (c *Credentials) Create(db *DB) error {
	return db.DB.Create(&c).Error
}

// Update will update an existing credentials entry in the database with the
// current object's values. It returns an error if the update fails.
func (c *Credentials) Update(db *DB) error {
	return db.DB.Model(&c).Where("id = ?", c.Id).Updates(&c).Error
}

// Delete will delete an existing credentials entry from the database that
// matches the current object's ID. It returns an error if the deletion fails.
func (c *Credentials) Delete(db *DB) error {
	return db.DB.Delete(&c).Error
}

// GetDefault will fetch the default credentials entry from the database and
// populate the current object with its values. It returns an error if the fetch
// fails.
func (c *Credentials) GetDefault(db *DB) error {
	return db.DB.Where("alias = ?", "Default").First(&c).Error
}

// GetById fetches a credentials entry from the database using the current object's ID
// and populates the current object with its values. It returns an error if the fetch fails.
func (c *Credentials) GetById(db *DB) error {
	return db.DB.First(&c, "id = ?", c.Id).Error
}

// GetAll retrieves all credentials entries from the database and returns them
// as a slice of Credentials pointers. It returns an error if the retrieval fails.
func (c *Credentials) GetAll(db *DB) ([]*Credentials, error) {
	var creds []*Credentials
	return creds, db.DB.Find(&creds).Error
}
