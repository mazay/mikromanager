package utils

import (
	"encoding/json"
	"errors"
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
	exists, _ := db.Exists("credentials", "alias", c.Alias)
	if exists {
		return errors.New(fmt.Sprintf("Alias '%s' already exists, please pick another name", c.Alias))
	}
	c.Created = time.Now()
	c.Updated = time.Now()
	inrec, _ := json.Marshal(c)
	json.Unmarshal(inrec, &inInterface)
	_, err := db.Insert("credentials", inInterface)
	return err
}

func (c *Credentials) Update(db *database.DB) error {
	var inInterface map[string]interface{}
	c.Updated = time.Now()
	inrec, _ := json.Marshal(c)
	json.Unmarshal(inrec, &inInterface)
	return db.Update("credentials", "_id", c.Id, inInterface)
}

func (c *Credentials) Delete(db *database.DB) error {
	return db.DeleteById("credentials", c.Id)
}

func (c *Credentials) GetDefault(db *database.DB) error {
	creds, err := db.FindByKeyValue("credentials", "alias", "Default")
	if err != nil {
		log.Fatal(err)
		return err
	}

	inrec, err := json.Marshal(creds)
	json.Unmarshal(inrec, c)
	return err
}

func (c *Credentials) GetById(db *database.DB) error {
	doc, err := db.FindById("credentials", c.Id)
	if err != nil {
		log.Fatal(err)
		return err
	}

	dj, err := json.Marshal(doc)
	json.Unmarshal(dj, c)

	return err
}

func (c *Credentials) GetAll(db *database.DB) ([]*Credentials, error) {
	var credList []*Credentials

	docs, err := db.FindAll("credentials")
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
