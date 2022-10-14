package utils

import (
	"encoding/json"
	"log"

	database "github.com/mazay/mikromanager/db"
)

type Credentials struct {
	Id                string `json:"_id"`
	Username          string `json:"username"`
	EncryptedPassword string `json:"encryptedPassword"`
	Default           bool   `json:"default"`
}

func (d *Credentials) GetAll(db *database.DB) ([]*Credentials, error) {
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
