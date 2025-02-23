package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSessionsCreate(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	user, err := createTestUser(db)
	if err != nil {
		t.Fatal(err)
	}

	session := &Session{
		UserId: user.Id,
	}

	err = session.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, session.Id)
	assert.NotEmpty(t, session.CreatedAt)
	assert.NotEmpty(t, session.UpdatedAt)
}

func TestSessionsGetById(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	user, err := createTestUser(db)
	if err != nil {
		t.Fatal(err)
	}

	session := &Session{
		UserId: user.Id,
	}

	err = session.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	fetchedSession := &Session{}
	fetchedSession.Id = session.Id
	err = fetchedSession.GetById(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, session.Id, fetchedSession.Id)
	assert.NotEmpty(t, fetchedSession.CreatedAt)
	assert.NotEmpty(t, fetchedSession.UpdatedAt)
}

func TestSessionsGetByUserId(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	user, err := createTestUser(db)
	if err != nil {
		t.Fatal(err)
	}

	session := &Session{
		UserId: user.Id,
	}

	err = session.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	fetchedSession := &Session{}
	fetchedSession.UserId = session.UserId
	sessionList, err := fetchedSession.GetByUserId(db, user.Id)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, sessionList, 1)
	assert.Equal(t, session.Id, sessionList[0].Id)
	assert.NotEmpty(t, sessionList[0].CreatedAt)
	assert.NotEmpty(t, sessionList[0].UpdatedAt)
}

func TestSessionsDelete(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	user, err := createTestUser(db)
	if err != nil {
		t.Fatal(err)
	}

	session := &Session{
		UserId: user.Id,
	}

	err = session.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	err = session.Delete(db)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSessionsExpired(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	user, err := createTestUser(db)
	if err != nil {
		t.Fatal(err)
	}

	session := &Session{
		UserId: user.Id,
		ValidThrough: time.Now().Add(
			-time.Hour * 2,
		),
	}

	err = session.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, session.Expired())
}
