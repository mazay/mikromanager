package db

import "gorm.io/gorm/clause"

type DeviceGroup struct {
	Base
	Name    string    `gorm:"unique"`
	Devices []*Device `gorm:"foreignKey:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (g *DeviceGroup) Create(db *DB) error {
	return db.DB.Omit(clause.Associations).Create(&g).Error
}

func (g *DeviceGroup) Update(db *DB) error {
	return db.DB.Omit(clause.Associations).Model(&g).Where("id = ?", g.Id).Updates(g).Error
}

func (g *DeviceGroup) Delete(db *DB) error {
	return db.DB.Omit(clause.Associations).Delete(&g).Error
}

func (g *DeviceGroup) GetAll(db *DB) ([]*DeviceGroup, error) {
	var groups []*DeviceGroup
	return groups, db.DB.Find(&groups).Error
}

func (g *DeviceGroup) GetById(db *DB) error {
	return db.DB.First(&g, "id = ?", g.Id).Error
}
