package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	database "github.com/mazay/mikromanager/db"
)

type Credentials struct {
	Id                string    `json:"_id"`
	Alias             string    `json:"alias"`
	Username          string    `json:"username"`
	EncryptedPassword string    `json:"encryptedPassword"`
	Created           time.Time `json:"created"`
	Updated           time.Time `json:"updated"`
}

func (c *Credentials) Create(db *database.DB) error {
	var inInterface map[string]interface{}
	// check if credentials with that alias already exist
	exists, _ := db.Exists(db.Collections["credentials"], "alias", c.Alias)
	if exists {
		return fmt.Errorf("Alias '%s' already exists, please pick another name", c.Alias)
	}
	c.Created = time.Now()
	c.Updated = time.Now()
	inrec, _ := json.Marshal(c)
	err := json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return err
	}
	_, err = db.Insert(db.Collections["credentials"], inInterface)
	return err
}

func (c *Credentials) Update(db *database.DB) error {
	var inInterface map[string]interface{}
	c.Updated = time.Now()
	inrec, _ := json.Marshal(c)
	err := json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return err
	}
	return db.Update(db.Collections["credentials"], "_id", c.Id, inInterface)
}

func (c *Credentials) Delete(db *database.DB) error {
	return db.DeleteById(db.Collections["credentials"], c.Id)
}

func (c *Credentials) GetDefault(db *database.DB) error {
	creds, err := db.FindByKeyValue(db.Collections["credentials"], "alias", "Default")
	if err != nil {
		log.Fatal(err)
		return err
	}

	inrec, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	err = json.Unmarshal(inrec, c)
	if err != nil {
		return err
	}
	return err
}

func (c *Credentials) GetById(db *database.DB) error {
	doc, err := db.FindById(db.Collections["credentials"], c.Id)
	if err != nil {
		log.Fatal(err)
		return err
	}

	dj, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dj, c)
	if err != nil {
		return err
	}

	return err
}

func (c *Credentials) GetAll(db *database.DB) ([]*Credentials, error) {
	var credList []*Credentials

	docs, err := db.FindAll(db.Collections["credentials"])
	if err != nil {
		log.Fatal(err)
		return credList, err
	}

	for _, doc := range docs {
		dm := &Credentials{}
		dj, _ := json.Marshal(doc)
		_ = json.Unmarshal(dj, dm)
		credList = append(credList, dm)
	}

	return credList, nil
}
