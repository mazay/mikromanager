package utils

import (
	"encoding/json"
)

type Credentials struct {
	Id                string `json:"_id"`
	Username          string `json:"username"`
	EncryptedPassword string `json:"encryptedPassword"`
	Default           bool   `json:"default"`
}

func (d *Credentials) FromListOfMaps(docs []map[string]string) []*Credentials {
	var credList []*Credentials

	for _, doc := range docs {
		dm := &Credentials{}
		dj, _ := json.Marshal(doc)
		_ = json.Unmarshal(dj, dm)
		credList = append(credList, dm)
	}

	return credList
}
