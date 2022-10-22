package utils

import (
	"encoding/json"
	"log"
	"os"
	"time"

	database "github.com/mazay/mikromanager/db"
)

type Export struct {
	Id       string    `json:"_id"`
	DeviceId string    `json:"deviceId"`
	Filename string    `json:"filename"`
	Created  time.Time `json:"created"`
}

func (b *Export) Create(db *database.DB) error {
	var inInterface map[string]interface{}
	b.Created = time.Now()
	inrec, _ := json.Marshal(b)
	err := json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return err
	}
	_, err = db.Insert(db.Collections["exports"], inInterface)
	return err
}

func (b *Export) Delete(db *database.DB) error {
	err := os.Remove(b.Filename)
	if err != nil {
		return err
	}
	return db.DeleteById(db.Collections["exports"], b.Id)
}

func (b *Export) GetById(db *database.DB) error {
	doc, err := db.FindById(db.Collections["exports"], b.Id)
	if err != nil {
		log.Fatal(err)
		return err
	}

	dj, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dj, b)
	if err != nil {
		return err
	}

	return err
}

func (b *Export) GetByDeviceId(db *database.DB, deviceId string) ([]*Export, error) {
	var exportList []*Export

	docs, err := db.FindAllByKeyValue(db.Collections["exports"], "deviceId", deviceId)
	if err != nil {
		log.Fatal(err)
		return exportList, err
	}

	for _, doc := range docs {
		bm := &Export{}
		dj, _ := json.Marshal(doc)
		_ = json.Unmarshal(dj, bm)
		exportList = append(exportList, bm)
	}

	return exportList, nil
}

func (b *Export) GetAll(db *database.DB) ([]*Export, error) {
	var exportList []*Export

	docs, err := db.FindAll(db.Collections["exports"])
	if err != nil {
		log.Fatal(err)
		return exportList, err
	}

	for _, doc := range docs {
		bm := &Export{}
		dj, _ := json.Marshal(doc)
		_ = json.Unmarshal(dj, bm)
		exportList = append(exportList, bm)
	}

	return exportList, nil
}
