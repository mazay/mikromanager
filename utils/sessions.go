package utils

import (
	"encoding/json"
	"fmt"
	"time"

	database "github.com/mazay/mikromanager/db"
)

type Session struct {
	Id           string    `json:"_id"`
	UserId       string    `json:"userId"`
	ValidThrough time.Time `json:"expires"`
}

func (s *Session) Expired() bool {
	if s.ValidThrough.IsZero() {
		return true
	}
	return s.ValidThrough.Before(time.Now())
}

func (s *Session) Create(db *database.DB) error {
	var inInterface map[string]interface{}
	s.ValidThrough = time.Now().Add(time.Hour * 24)
	fmt.Printf("%s", s.ValidThrough)
	inrec, _ := json.Marshal(s)
	err := json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return err
	}
	s.Id, err = db.Insert(db.Collections["sessions"], inInterface)
	return err
}

func (s *Session) Delete(db *database.DB) error {
	return db.DeleteById(db.Collections["sessions"], s.Id)
}

func (s *Session) GetById(db *database.DB) error {
	doc, err := db.FindById(db.Collections["sessions"], s.Id)
	if err != nil {
		return err
	}

	dj, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dj, s)
	if err != nil {
		return err
	}

	return err
}

func (s *Session) Update(db *database.DB) error {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(s)
	err := json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return err
	}
	return db.Update(db.Collections["sessions"], "_id", s.Id, inInterface)
}

func (s *Session) GetAll(db *database.DB) ([]*User, error) {
	var userList []*User

	docs, err := db.FindAll(db.Collections["sessions"])
	if err != nil {
		return userList, err
	}

	for _, doc := range docs {
		um := &User{}
		uj, _ := json.Marshal(doc)
		_ = json.Unmarshal(uj, um)
		userList = append(userList, um)
	}

	return userList, nil
}

func (s *Session) GetByUserId(db *database.DB, userId string) ([]*Session, error) {
	var sessionList []*Session

	db.Sort("expires", -1)
	docs, err := db.FindAllByKeyValue(db.Collections["sessions"], "userId", userId)
	if err != nil {
		return sessionList, err
	}

	for _, doc := range docs {
		sm := &Session{}
		sj, _ := json.Marshal(doc)
		_ = json.Unmarshal(sj, sm)
		sessionList = append(sessionList, sm)
	}

	return sessionList, nil
}
