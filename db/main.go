package db

import (
	"log"

	"github.com/ostafen/clover"
)

var (
	collections = []string{
		"credentials",
		"devices",
	}
)

type DB struct {
	Path string
	api  *clover.DB
}

func (db *DB) Init() (bool, error) {
	// open or create the DB
	database, err := clover.Open(db.Path)
	if err != nil {
		return false, err
	}
	db.api = database
	// create the set of collections
	for _, collection := range collections {
		if !db.HasCollection(collection) {
			db.CreateCollection(collection)
		}
	}

	return true, nil
}

func (db *DB) Close() error {
	return db.api.Close()
}

func (db *DB) HasCollection(collection string) bool {
	exists, err := db.api.HasCollection(collection)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

func (db *DB) CreateCollection(collection string) error {
	return db.api.CreateCollection(collection)
}

func (db *DB) Insert(collection string, values map[string]interface{}) (string, error) {
	doc := clover.NewDocument()
	doc.SetAll(values)

	return db.api.InsertOne(collection, doc)
}

func (db *DB) Update(collection string, key string, value string, updates map[string]interface{}) error {
	exists, _ := db.Exists(collection, key, value)
	if exists {
		return db.api.Query(collection).Where(clover.Field(key).Eq(value)).Update(updates)
	} else {
		updates[key] = value
		_, err := db.Insert(collection, updates)
		return err
	}
}

func (db *DB) DeleteById(collection string, id string) error {
	return db.api.Query(collection).DeleteById(id)
}

func (db *DB) Exists(collection string, key string, value string) (bool, error) {
	return db.api.Query(collection).Where(clover.Field(key).Eq(value)).Exists()
}

func (db *DB) FindById(collection string, id string) (map[string]string, error) {
	var inInterface map[string]string
	query := db.api.Query(collection)
	doc, err := query.FindById(id)
	doc.Unmarshal(&inInterface)
	return inInterface, err
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
