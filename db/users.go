package db

type User struct {
	Base
	Username          string `gorm:"unique"`
	EncryptedPassword string
}

func (u *User) Create(db *DB) error {
	return db.DB.Create(&u).Error
}

func (u *User) Delete(db *DB) error {
	return db.DB.Delete(&u).Error
}

func (u *User) GetById(db *DB) error {
	return db.DB.First(&u, "id = ?", u.Id).Error
}

func (u *User) GetByUsername(db *DB) error {
	return db.DB.First(&u, "username = ?", u.Username).Error
}

func (u *User) Update(db *DB) error {
	return db.DB.Model(&u).Where("id = ?", u.Id).Updates(u).Error
}

func (u *User) GetAll(db *DB) ([]*User, error) {
	var userList []*User
	return userList, db.DB.Find(&userList).Error
}
