package db

type DeviceGroup struct {
	Base
	Name    string    `gorm:"unique"`
	Devices []*Device `gorm:"many2many:device_groups_devices;"`
}

// Create will create a new device group entry in the database with the current object's values.
// The function returns an error if the creation fails.
func (g *DeviceGroup) Create(db *DB) error {
	return db.DB.Create(&g).Error
}

// Update will update an existing device group entry in the database with the current
// object's values, including its associated devices. If the update fails, it
// returns an error.
func (g *DeviceGroup) Update(db *DB) error {
	err := db.DB.Model(&g).Association("Devices").Replace(g.Devices)
	if err != nil {
		return err
	}
	return db.DB.Save(&g).Error
}

// Save will persist the current state of the device group to the database.
// If the save operation fails, it returns an error.
func (g *DeviceGroup) Save(db *DB) error {
	return db.DB.Save(&g).Error
}

// Delete will delete an existing device group entry from the database that matches the
// current object's ID. It returns an error if the deletion fails.
func (g *DeviceGroup) Delete(db *DB) error {
	return db.DB.Delete(&g).Error
}

// GetAllPlain retrieves all device group entries from the database and returns them
// as a slice of *DeviceGroup instances. It does not preload any associated devices.
// The function returns an error if the retrieval fails.
func (g *DeviceGroup) GetAllPlain(db *DB) ([]*DeviceGroup, error) {
	var groups []*DeviceGroup
	return groups, db.DB.Find(&groups).Error
}

// GetAllPreload retrieves all device group entries from the database and returns them
// as a slice of *DeviceGroup instances. It preloads the associated devices for each
// device group. The function returns an error if the retrieval fails.
func (g *DeviceGroup) GetAllPreload(db *DB) ([]*DeviceGroup, error) {
	var groups []*DeviceGroup
	return groups, db.DB.Model(&g).Preload("Devices").Find(&groups).Error
}

// GetById fetches a device group entry from the database using the current object's ID
// and populates the current object with its values, including its associated devices.
// It returns an error if the fetch fails.
func (g *DeviceGroup) GetById(db *DB) error {
	return db.DB.Model(g).Preload("Devices").First(&g, "id = ?", g.Id).Error
}

// Append will append the given devices to the current device group in the database.
// The function returns an error if the append fails.
func (g *DeviceGroup) Append(db *DB, devices []*Device) error {
	return db.DB.Model(g).Association("Devices").Append(devices)
}

// Remove will remove the given devices from the current device group in the database.
// The function returns an error if the removal fails.
func (g *DeviceGroup) Remove(db *DB, devices []*Device) error {
	return db.DB.Model(g).Association("Devices").Delete(devices)
}
