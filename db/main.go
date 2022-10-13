package db

import (
	"log"

	"github.com/ostafen/clover"
)

type DB struct {
	Path string
	api  *clover.DB
}

func (db *DB) Init() (bool, error) {
	database, err := clover.Open(db.Path)
	if err != nil {
		return false, err
	}

	db.api = database

	exists, _ := db.api.HasCollection("devices")
	if !exists {
		db.api.CreateCollection("devices")
	}

	return true, nil
}

func (db *DB) Close() error {
	return db.api.Close()
}

func (db *DB) Insert(collection string, values map[string]interface{}) (string, error) {
	doc := clover.NewDocument()
	doc.SetAll(values)

	return db.api.InsertOne(collection, doc)
}

func (db *DB) Update(collection string, key string, value string, updates map[string]interface{}) error {
	exists, _ := db.api.Query(collection).Where(clover.Field(key).Eq(value)).Exists()
	if exists {
		return db.api.Query(collection).Where(clover.Field(key).Eq(value)).Update(updates)
	} else {
		updates[key] = value
		_, err := db.Insert(collection, updates)
		return err
	}
}

func (db *DB) FindAll(collection string) ([]map[string]string, error) {
	var result []map[string]string
	query := db.api.Query(collection)
	docs, err := query.FindAll()
	for _, doc := range docs {
		var inInterface map[string]string
		doc.Unmarshal(&inInterface)
		result = append(result, inInterface)
	}
	return result, err
}

func (db *DB) Print() {
	query := db.api.Query("devices")
	docs, _ := query.FindAll()

	for _, doc := range docs {
		log.Println(doc)
	}
}

func (db *DB) Export(collection string, filename string) {
	db.api.ExportCollection(collection, filename)
}
