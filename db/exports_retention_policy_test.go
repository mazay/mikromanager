package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testExportsRetentionPolicy = &ExportsRetentionPolicy{
		Name:   "Default",
		Hourly: 24,
		Daily:  14,
		Weekly: 26,
	}
)

func TestExportsRetentionPolicyCreate(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	err = testExportsRetentionPolicy.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Default", testExportsRetentionPolicy.Name)
	assert.Equal(t, int64(24), testExportsRetentionPolicy.Hourly)
	assert.Equal(t, int64(14), testExportsRetentionPolicy.Daily)
	assert.Equal(t, int64(26), testExportsRetentionPolicy.Weekly)
	assert.NotEmpty(t, testExportsRetentionPolicy.Id)
	assert.NotEmpty(t, testExportsRetentionPolicy.CreatedAt)
	assert.NotEmpty(t, testExportsRetentionPolicy.UpdatedAt)
}

func TestExportsRetentionPolicyGetDefault(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	err = testExportsRetentionPolicy.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	err = testExportsRetentionPolicy.GetDefault(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Default", testExportsRetentionPolicy.Name)
	assert.Equal(t, int64(24), testExportsRetentionPolicy.Hourly)
	assert.Equal(t, int64(14), testExportsRetentionPolicy.Daily)
	assert.Equal(t, int64(26), testExportsRetentionPolicy.Weekly)
	assert.NotEmpty(t, testExportsRetentionPolicy.Id)
	assert.NotEmpty(t, testExportsRetentionPolicy.CreatedAt)
	assert.NotEmpty(t, testExportsRetentionPolicy.UpdatedAt)
}

// the update test should go last because it updates the default values
func TestExportsRetentionPolicyUpdate(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	err = testExportsRetentionPolicy.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	testExportsRetentionPolicy.Hourly = 25
	err = testExportsRetentionPolicy.Update(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Default", testExportsRetentionPolicy.Name)
	assert.Equal(t, int64(25), testExportsRetentionPolicy.Hourly)
	assert.Equal(t, int64(14), testExportsRetentionPolicy.Daily)
	assert.Equal(t, int64(26), testExportsRetentionPolicy.Weekly)
	assert.NotEmpty(t, testExportsRetentionPolicy.Id)
	assert.NotEmpty(t, testExportsRetentionPolicy.CreatedAt)
	assert.NotEmpty(t, testExportsRetentionPolicy.UpdatedAt)
}
