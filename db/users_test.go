package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testUser = &User{
		Username:          "test-user",
		EncryptedPassword: "test-password",
	}
)

func TestUserCreate(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	err = testUser.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "test-user", testUser.Username)
	assert.NotEmpty(t, testUser.Id)
	assert.NotEmpty(t, testUser.CreatedAt)
	assert.NotEmpty(t, testUser.UpdatedAt)
}

func TestUserGetByUsername(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	err = testUser.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	fetchedUser := &User{}
	fetchedUser.Username = testUser.Username
	err = fetchedUser.GetByUsername(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testUser.Id, fetchedUser.Id)
	assert.Equal(t, testUser.Username, fetchedUser.Username)
	assert.NotEmpty(t, fetchedUser.Id)
	assert.NotEmpty(t, fetchedUser.CreatedAt)
	assert.NotEmpty(t, fetchedUser.UpdatedAt)
}

func TestUserGetById(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	err = testUser.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	fetchedUser := &User{}
	fetchedUser.Id = testUser.Id
	err = fetchedUser.GetById(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testUser.Id, fetchedUser.Id)
	assert.Equal(t, testUser.Username, fetchedUser.Username)
	assert.NotEmpty(t, fetchedUser.Id)
	assert.NotEmpty(t, fetchedUser.CreatedAt)
	assert.NotEmpty(t, fetchedUser.UpdatedAt)
}

func TestUserUniqueUsername(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	err = testUser.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	nonUniqueUser := &User{
		Username:          "test-user",
		EncryptedPassword: "test-password",
	}

	err = nonUniqueUser.Create(db)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUserGetAll(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	err = testUser.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	userList, err := testUser.GetAll(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, userList, 1)
	assert.Equal(t, testUser.Id, userList[0].Id)
	assert.Equal(t, testUser.Username, userList[0].Username)
	assert.NotEmpty(t, userList[0].Id)
	assert.NotEmpty(t, userList[0].CreatedAt)
	assert.NotEmpty(t, userList[0].UpdatedAt)
}

func TestUserDelete(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	err = testUser.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	err = testUser.Delete(db)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserUpdate(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	err = testUser.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	testUser.Username = "test-user-updated"
	err = testUser.Update(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "test-user-updated", testUser.Username)
	assert.NotEmpty(t, testUser.Id)
	assert.NotEmpty(t, testUser.CreatedAt)
	assert.NotEmpty(t, testUser.UpdatedAt)
}
