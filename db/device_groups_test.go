package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeviceGroupCreate(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	dev, err := createTestDevice(db)
	if err != nil {
		t.Fatal(err)
	}

	group := &DeviceGroup{
		Name: "test-group",
		Devices: []*Device{
			dev,
		},
	}

	err = group.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, group.Devices, 1)
	assert.NotEmpty(t, group.Id)
	assert.NotEmpty(t, group.CreatedAt)
	assert.NotEmpty(t, group.UpdatedAt)
}

func TestDeviceGroupUpdate(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	dev, err := createTestDevice(db)
	if err != nil {
		t.Fatal(err)
	}

	group := &DeviceGroup{
		Name: "test-group",
		Devices: []*Device{
			dev,
		},
	}

	err = group.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	group.Name = "test-group-updated"
	err = group.Update(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "test-group-updated", group.Name)
	assert.Len(t, group.Devices, 1)
	assert.NotEmpty(t, group.Id)
	assert.NotEmpty(t, group.CreatedAt)
	assert.NotEmpty(t, group.UpdatedAt)
}

func TestDeviceGroupDelete(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	group := &DeviceGroup{
		Name: "test-group",
	}

	err = group.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	err = group.Delete(db)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeviceGroupGetAllPreload(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	dev, err := createTestDevice(db)
	if err != nil {
		t.Fatal(err)
	}

	group := &DeviceGroup{
		Name: "test-group",
		Devices: []*Device{
			dev,
		},
	}

	err = group.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	groups, err := group.GetAllPreload(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, groups)
	assert.Len(t, groups, 1)
	assert.Len(t, groups[0].Devices, 1)
	assert.Equal(t, group.Id, groups[0].Id)
	assert.Equal(t, group.Name, groups[0].Name)
	assert.NotEmpty(t, groups[0].Id)
	assert.NotEmpty(t, groups[0].CreatedAt)
	assert.NotEmpty(t, groups[0].UpdatedAt)
}

func TestDeviceGroupGetAllPlain(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	dev, err := createTestDevice(db)
	if err != nil {
		t.Fatal(err)
	}

	group := &DeviceGroup{
		Name: "test-group",
		Devices: []*Device{
			dev,
		},
	}

	err = group.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	groups, err := group.GetAllPlain(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, groups)
	assert.Len(t, groups, 1)
	assert.Len(t, groups[0].Devices, 0)
	assert.Equal(t, group.Id, groups[0].Id)
	assert.Equal(t, group.Name, groups[0].Name)
	assert.NotEmpty(t, groups[0].Id)
	assert.NotEmpty(t, groups[0].CreatedAt)
	assert.NotEmpty(t, groups[0].UpdatedAt)
}

func TestDeviceGroupGetById(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	group := &DeviceGroup{
		Name: "test-group",
	}

	err = group.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	fetchedGroup := &DeviceGroup{}
	fetchedGroup.Id = group.Id
	err = fetchedGroup.GetById(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, group.Id, fetchedGroup.Id)
	assert.Equal(t, group.Name, fetchedGroup.Name)
	assert.NotEmpty(t, fetchedGroup.Id)
	assert.NotEmpty(t, fetchedGroup.CreatedAt)
	assert.NotEmpty(t, fetchedGroup.UpdatedAt)
}
