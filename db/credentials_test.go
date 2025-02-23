package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredentialsCreate(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	creds := &Credentials{
		Alias:             "test-alias",
		Username:          "test-username",
		EncryptedPassword: "test-password",
	}

	err = creds.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "test-alias", creds.Alias)
	assert.Equal(t, "test-username", creds.Username)
	assert.Equal(t, "test-password", creds.EncryptedPassword)
	assert.NotEmpty(t, creds.Id)
	assert.NotEmpty(t, creds.CreatedAt)
	assert.NotEmpty(t, creds.UpdatedAt)
	assert.NotEmpty(t, creds.EncryptedPassword)

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCredentialsGet(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	creds := &Credentials{
		Alias:             "test-alias",
		Username:          "test-username",
		EncryptedPassword: "test-password",
	}

	err = creds.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	fetchedCreds := &Credentials{}
	fetchedCreds.Id = creds.Id
	err = fetchedCreds.GetById(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "test-alias", fetchedCreds.Alias)
	assert.Equal(t, "test-username", fetchedCreds.Username)
	assert.Equal(t, "test-password", fetchedCreds.EncryptedPassword)
	assert.NotEmpty(t, fetchedCreds.Id)
	assert.NotEmpty(t, fetchedCreds.CreatedAt)
	assert.NotEmpty(t, fetchedCreds.UpdatedAt)
	assert.NotEmpty(t, fetchedCreds.EncryptedPassword)

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCredentialsUpdate(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	creds := &Credentials{
		Alias:             "test-alias",
		Username:          "test-username",
		EncryptedPassword: "test-password",
	}

	err = creds.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	creds.Alias = "updated-alias"
	creds.Username = "updated-username"
	creds.EncryptedPassword = "updated-password"

	err = creds.Update(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "updated-alias", creds.Alias)
	assert.Equal(t, "updated-username", creds.Username)
	assert.Equal(t, "updated-password", creds.EncryptedPassword)

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCredentialsDelete(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	creds := &Credentials{
		Alias:             "test-alias",
		Username:          "test-username",
		EncryptedPassword: "test-password",
	}

	err = creds.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	err = creds.Delete(db)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCredentialsGetDefault(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	creds := &Credentials{
		Alias:             "Default",
		Username:          "test-username",
		EncryptedPassword: "test-password",
	}

	err = creds.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	fetchedCreds := &Credentials{}
	err = fetchedCreds.GetDefault(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Default", fetchedCreds.Alias)
	assert.Equal(t, "test-username", fetchedCreds.Username)
	assert.Equal(t, "test-password", fetchedCreds.EncryptedPassword)
	assert.NotEmpty(t, fetchedCreds.Id)
	assert.NotEmpty(t, fetchedCreds.CreatedAt)
	assert.NotEmpty(t, fetchedCreds.UpdatedAt)
	assert.NotEmpty(t, fetchedCreds.EncryptedPassword)

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCredentialsGetAll(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	creds := &Credentials{
		Alias:             "test-alias",
		Username:          "test-username",
		EncryptedPassword: "test-password",
	}

	err = creds.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	fetchedCreds, err := creds.GetAll(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "test-alias", fetchedCreds[0].Alias)
	assert.Equal(t, "test-username", fetchedCreds[0].Username)
	assert.Equal(t, "test-password", fetchedCreds[0].EncryptedPassword)
	assert.NotEmpty(t, fetchedCreds[0].Id)
	assert.NotEmpty(t, fetchedCreds[0].CreatedAt)
	assert.NotEmpty(t, fetchedCreds[0].UpdatedAt)
	assert.NotEmpty(t, fetchedCreds[0].EncryptedPassword)

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
}
