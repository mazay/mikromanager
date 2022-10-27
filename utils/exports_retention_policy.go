package utils

import (
	"encoding/json"
	"fmt"

	database "github.com/mazay/mikromanager/db"
)

type ExportsRetentionPolicy struct {
	Id     string `json:"_id"`
	Name   string `json:"name"`
	Hourly int64  `json:"hourly"`
	Daily  int64  `json:"daily"`
	Weekly int64  `json:"weekly"`
}

func (rp *ExportsRetentionPolicy) Create(db *database.DB) error {
	var inInterface map[string]interface{}
	// check if credentials with that alias already exist
	exists, _ := db.Exists(db.Collections["exportsRetentionPolicy"], "name", rp.Name)
	if exists {
		return fmt.Errorf("retention poilicy '%s' already exists, please pick another name", rp.Name)
	}
	inrec, _ := json.Marshal(rp)
	err := json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return err
	}
	_, err = db.Insert(db.Collections["exportsRetentionPolicy"], inInterface)
	return err
}

func (rp *ExportsRetentionPolicy) Update(db *database.DB) error {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(rp)
	err := json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return err
	}
	return db.Update(db.Collections["exportsRetentionPolicy"], "_id", rp.Id, inInterface)
}

func (rp *ExportsRetentionPolicy) GetDefault(db *database.DB) error {
	policy, err := db.FindByKeyValue(db.Collections["exportsRetentionPolicy"], "name", "Default")
	if err != nil {
		return err
	}

	inrec, err := json.Marshal(policy)
	if err != nil {
		return err
	}
	err = json.Unmarshal(inrec, rp)
	if err != nil {
		return err
	}
	return err
}
