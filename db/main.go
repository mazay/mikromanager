package db

import (
	"github.com/ostafen/clover"
	"go.uber.org/zap"
)

type DB struct {
	Path        string
	api         *clover.DB
	Logger      *zap.Logger
	Collections map[string]string
	SortBy      *clover.SortOption
}

func (db *DB) Init() error {
	// open or create the DB
	database, err := clover.Open(db.Path)
	if err != nil {
		return err
	}
	db.api = database
	db.Collections = collectionsMap()
	// create the set of collections
	for _, name := range db.Collections {
		if !db.HasCollection(db.Collections[name]) {
			if db.CreateCollection(db.Collections[name]) != nil {
				return err
			}
		}
	}

	return nil
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
		db.Logger.Error(err.Error())
	}
	return exists
}

func (db *DB) CreateCollection(collection string) error {
	return db.api.CreateCollection(collection)
}

func (db *DB) Sort(field string, direction int) {
	db.SortBy = &clover.SortOption{
		Field:     field,
		Direction: direction,
	}
}

func (db *DB) getBaseQuery(collection string) *clover.Query {
	baseQuery := db.api.Query(collection)
	if db.SortBy != (&clover.SortOption{}) {
		baseQuery = db.api.Query(collection).Sort(*db.SortBy)
	}
	return baseQuery
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

func (db *DB) FindById(collection string, id string) (map[string]interface{}, error) {
	var inInterface map[string]interface{}
	query := db.api.Query(collection)
	doc, err := query.FindById(id)
	if err == nil && doc != nil {
		err = doc.Unmarshal(&inInterface)
	}
	return inInterface, err
}

func (db *DB) FindByKeyValue(collection string, key string, value string) (map[string]interface{}, error) {
	var inInterface map[string]interface{}
	doc, err := db.api.Query(collection).Where(clover.Field(key).Eq(value)).FindFirst()
	if err == nil && doc != nil {
		err = doc.Unmarshal(&inInterface)
	}
	return inInterface, err
}

func (db *DB) FindAllByKeyValue(collection string, key string, value string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	baseQuery := db.getBaseQuery(collection)
	docs, err := baseQuery.Where(clover.Field(key).Eq(value)).FindAll()
	for _, doc := range docs {
		var inInterface map[string]interface{}
		if doc.Unmarshal(&inInterface) == nil {
			result = append(result, inInterface)
		}
	}
	return result, err
}

func (db *DB) FindByCriteriaSet(collection string, set []map[string]string) (map[string]interface{}, error) {
	var inInterface map[string]interface{}
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
	baseQuery := db.getBaseQuery(collection)
	doc, err := baseQuery.Where(criteria).FindFirst()
	if err == nil && doc != nil {
		err = doc.Unmarshal(&inInterface)
	}
	return inInterface, err
}

func (db *DB) FindAll(collection string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	baseQuery := db.getBaseQuery(collection)
	docs, err := baseQuery.FindAll()
	for _, doc := range docs {
		var inInterface map[string]interface{}
		if doc.Unmarshal(&inInterface) == nil {
			result = append(result, inInterface)
		}
	}
	return result, err
}

func (db *DB) Export(collection string, filename string) error {
	return db.api.ExportCollection(collection, filename)
}
