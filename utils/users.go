package utils

import (
	"encoding/json"
	"fmt"
	"time"

	database "github.com/mazay/mikromanager/db"
)

type User struct {
	Id                string    `json:"_id"`
	Username          string    `json:"username"`
	EncryptedPassword string    `json:"encryptedPassword"`
	Created           time.Time `json:"created"`
	Updated           time.Time `json:"updated"`
}

func (u *User) Create(db *database.DB) error {
	var inInterface map[string]interface{}
	// check if username already taken
	exists, _ := db.Exists(db.Collections["users"], "username", u.Username)
	if exists {
		return fmt.Errorf("User '%s' already exists", u.Username)
	}
	u.Created = time.Now()
	inrec, _ := json.Marshal(u)
	err := json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return err
	}
	_, err = db.Insert(db.Collections["users"], inInterface)
	return err
}

func (u *User) Delete(db *database.DB) error {
	var (
		err     error
		session = &Session{}
	)

	// delete sessions first
	sessions, err := session.GetByUserId(db, u.Id)
	if err != nil {
		return err
	}

	for _, s := range sessions {
		err = s.Delete(db)
		if err != nil {
			return err
		}
	}

	return db.DeleteById(db.Collections["users"], u.Id)
}

func (u *User) GetById(db *database.DB) error {
	doc, err := db.FindById(db.Collections["users"], u.Id)
	if err != nil {
		return err
	}

	dj, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dj, u)
	if err != nil {
		return err
	}

	return err
}

func (u *User) GetByUsername(db *database.DB) error {
	doc, err := db.FindByKeyValue(db.Collections["users"], "username", u.Username)
	if err != nil {
		return err
	}

	dj, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dj, u)
	if err != nil {
		return err
	}

	return err
}

func (u *User) Update(db *database.DB) error {
	var inInterface map[string]interface{}
	u.Updated = time.Now()
	inrec, _ := json.Marshal(u)
	err := json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return err
	}
	return db.Update(db.Collections["users"], "_id", u.Id, inInterface)
}

func (u *User) GetAll(db *database.DB) ([]*User, error) {
	var userList []*User

	db.Sort("username", 1)
	docs, err := db.FindAll(db.Collections["users"])
	if err != nil {
		return userList, err
	}

	for _, doc := range docs {
		um := &User{}
		uj, _ := json.Marshal(doc)
		_ = json.Unmarshal(uj, um)
		userList = append(userList, um)
	}

	return userList, nil
}
