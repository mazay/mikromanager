package db

type ConfigurationSnippet struct {
	Base
	Name    string `gorm:"unique"`
	Content string
	Groups  []*DeviceGroup `gorm:"many2many:device_groups_configuration_snippets;"`
}

func (s *ConfigurationSnippet) Create(db *DB) error {
	return db.DB.Create(&s).Error
}

func (s *ConfigurationSnippet) Update(db *DB) error {
	return db.DB.Save(&s).Error
}

func (s *ConfigurationSnippet) Delete(db *DB) error {
	return db.DB.Delete(&s).Error
}

func (s *ConfigurationSnippet) GetAll(db *DB) ([]*ConfigurationSnippet, error) {
	var snippets []*ConfigurationSnippet
	if err := db.DB.Find(&snippets).Error; err != nil {
		return nil, err
	}
	return snippets, nil
}

func (s *ConfigurationSnippet) GetById(db *DB) error {
	return db.DB.First(&s).Error
}
