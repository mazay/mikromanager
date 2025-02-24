package db

type DeviceGroup struct {
	Base
	Name    string    `gorm:"unique"`
	Devices []*Device `gorm:"many2many:device_groups_devices;"`
}

func (g *DeviceGroup) Create(db *DB) error {
	return db.DB.Create(&g).Error
}

func (g *DeviceGroup) Update(db *DB) error {
	err := db.DB.Model(&g).Association("Devices").Replace(g.Devices)
	if err != nil {
		return err
	}
	return db.DB.Save(&g).Error
}

func (g *DeviceGroup) Delete(db *DB) error {
	return db.DB.Delete(&g).Error
}

func (g *DeviceGroup) GetAll(db *DB) ([]*DeviceGroup, error) {
	var groups []*DeviceGroup
	return groups, db.DB.Model(g).Preload("Devices").Find(&groups).Error
}

func (g *DeviceGroup) GetById(db *DB) error {
	return db.DB.Model(g).Preload("Devices").First(&g, "id = ?", g.Id).Error
}
