package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	database "github.com/mazay/mikromanager/db"
)

type Credentials struct {
	Id                string `json:"_id"`
	Alias             string `json:"alias"`
	Username          string `json:"username"`
	EncryptedPassword string `json:"encryptedPassword"`
}

func (c *Credentials) Create(db *database.DB) error {
	var inInterface map[string]interface{}
	// check if credentials with that alias already exist
	exists, _ := db.Exists("credentials", "alias", c.Alias)
	if exists {
		return errors.New(fmt.Sprintf("Alias '%s' already exists, please pick another name", c.Alias))
	}
	inrec, _ := json.Marshal(c)
	json.Unmarshal(inrec, &inInterface)
	_, err := db.Insert("credentials", inInterface)
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
