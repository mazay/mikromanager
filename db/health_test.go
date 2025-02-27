package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthSave(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	device, err := createTestDevice(db)
	if err != nil {
		t.Fatal(err)
	}

	health := &Health{
		DeviceId: device.Id,
		Voltage:  3.3,
	}

	assert.NoError(t, health.Save(db))
	assert.Equal(t, device.Id, health.Id)
	assert.Equal(t, device.Id, health.DeviceId)
	assert.Equal(t, health.Voltage, float32(3.3))
}

func TestHealthDelete(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	device, err := createTestDevice(db)
	if err != nil {
		t.Fatal(err)
	}

	health := &Health{
		DeviceId: device.Id,
		Voltage:  3.3,
	}

	assert.NoError(t, health.Save(db))
	assert.NoError(t, health.Delete(db))
}

func TestHealthGetByDeviceId(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	device, err := createTestDevice(db)
	if err != nil {
		t.Fatal(err)
	}

	health := &Health{
		DeviceId: device.Id,
		Voltage:  3.3,
	}

	helthFetched := &Health{DeviceId: device.Id}

	assert.NoError(t, health.Save(db))
	assert.NoError(t, helthFetched.GetByDeviceId(db))
	assert.Equal(t, health.Id, helthFetched.Id)
	assert.Equal(t, health.DeviceId, helthFetched.DeviceId)
	assert.Equal(t, health.Voltage, helthFetched.Voltage)
}
