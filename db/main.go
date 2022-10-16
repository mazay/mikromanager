package db

import (
	"github.com/ostafen/clover"
	"github.com/sirupsen/logrus"
)

type DB struct {
	Path        string
	api         *clover.DB
	Logger      *logrus.Entry
	Collections map[string]string
}

func (db *DB) Init() (bool, error) {
	// open or create the DB
	database, err := clover.Open(db.Path)
	if err != nil {
		return false, err
	}
	db.api = database
	db.Collections = collectionsMap()
	// create the set of collections
	for _, name := range db.Collections {
		if !db.HasCollection(db.Collections[name]) {
			db.CreateCollection(db.Collections[name])
		}
	}

	return true, nil
}

func (db *DB) Close() error {
	return db.api.Close()
}

func (db *DB) ListCollections() ([]string, error) {
	return db.api.ListCollections()
}

func (db *DB) HasCollection(collection string) bool {
	exists, err := db.api.HasCollection(collection)
	if err != nil {
		db.Logger.Error(err)
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

func (db *DB) UpdateById(collection string, id string, updates map[string]interface{}) error {
	return db.api.Query(collection).UpdateById(id, updates)
}

func (db *DB) DeleteById(collection string, id string) error {
	return db.api.Query(collection).DeleteById(id)
}

func (db *DB) Exists(collection string, key string, value string) (bool, error) {
	return db.api.Query(collection).Where(clover.Field(key).Eq(value)).Exists()
}

func (db *DB) ExistsByCriteriaSet(collection string, set []map[string]string) (bool, error) {
	var criteria = &clover.Criteria{}
	for idx, crt := range set {
		for k, v := range crt {
			if idx == 0 {
				criteria = clover.Field(k).Eq(v)
			} else {
				criteria = criteria.And(clover.Field(k).Eq(v))
			}
		}
	}
	return db.api.Query(collection).Where(criteria).Exists()
}

func (db *DB) FindById(collection string, id string) (map[string]string, error) {
	var inInterface map[string]string
	query := db.api.Query(collection)
	doc, err := query.FindById(id)
	if err == nil && doc != nil {
		err = doc.Unmarshal(&inInterface)
	}
	return inInterface, err
}

func (db *DB) FindByKeyValue(collection string, key string, value string) (map[string]string, error) {
	var inInterface map[string]string
	doc, err := db.api.Query(collection).Where(clover.Field(key).Eq(value)).FindFirst()
	if err == nil && doc != nil {
		err = doc.Unmarshal(&inInterface)
	}
	return inInterface, err
}

func (db *DB) FindByCriteriaSet(collection string, set []map[string]string) (map[string]string, error) {
	var inInterface map[string]string
	var criteria = &clover.Criteria{}
	for idx, crt := range set {
		for k, v := range crt {
			if idx == 0 {
				criteria = clover.Field(k).Eq(v)
			} else {
				criteria = criteria.And(clover.Field(k).Eq(v))
			}
		}
	}
	doc, err := db.api.Query(collection).Where(criteria).FindFirst()
	if err == nil && doc != nil {
		err = doc.Unmarshal(&inInterface)
	}
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

func (db *DB) Export(collection string, filename string) {
	db.api.ExportCollection(collection, filename)
}
