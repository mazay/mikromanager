package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testExport = &Export{
		S3Key:        "testKey",
		LastModified: nil,
		ETag:         "testETag",
		Size:         nil,
	}
)

func TestExportSave(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	device, err := createTestDevice(db)
	if err != nil {
		t.Fatal(err)
	}

	export := testExport
	export.Device = device
	err = export.Save(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, export.Id)
	assert.NotEmpty(t, export.CreatedAt)
	assert.NotEmpty(t, export.UpdatedAt)
	assert.NotEmpty(t, export.DeviceId)
}

func TestExportDelete(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	device, err := createTestDevice(db)
	if err != nil {
		t.Fatal(err)
	}

	export := testExport
	export.Device = device
	err = export.Save(db)
	if err != nil {
		t.Fatal(err)
	}

	err = export.Delete(db)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExportGetById(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	device, err := createTestDevice(db)
	if err != nil {
		t.Fatal(err)
	}

	export := testExport
	export.Device = device
	err = export.Save(db)
	if err != nil {
		t.Fatal(err)
	}

	fetchedExport := &Export{}
	fetchedExport.Id = export.Id
	err = fetchedExport.GetById(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, export.Id, fetchedExport.Id)
	assert.Equal(t, export.S3Key, fetchedExport.S3Key)
	assert.Equal(t, export.ETag, fetchedExport.ETag)
	assert.Equal(t, export.Size, fetchedExport.Size)
	assert.Equal(t, export.DeviceId, fetchedExport.DeviceId)
	assert.NotEmpty(t, export.CreatedAt)
	assert.NotEmpty(t, export.UpdatedAt)
}

func TestExportGetByDeviceId(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	device, err := createTestDevice(db)
	if err != nil {
		t.Fatal(err)
	}

	export := testExport
	export.Device = device
	err = export.Save(db)
	if err != nil {
		t.Fatal(err)
	}

	fetchedExports, err := export.GetByDeviceId(db, export.DeviceId)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(fetchedExports))
	assert.Equal(t, export.Id, fetchedExports[0].Id)
	assert.Equal(t, export.S3Key, fetchedExports[0].S3Key)
	assert.Equal(t, export.ETag, fetchedExports[0].ETag)
	assert.Equal(t, export.Size, fetchedExports[0].Size)
	assert.Equal(t, export.DeviceId, fetchedExports[0].DeviceId)
	assert.NotEmpty(t, export.CreatedAt)
	assert.NotEmpty(t, export.UpdatedAt)
}

func TestExportGetAll(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	device, err := createTestDevice(db)
	if err != nil {
		t.Fatal(err)
	}

	export := testExport
	export.Device = device
	err = export.Save(db)
	if err != nil {
		t.Fatal(err)
	}

	fetchedExports, err := export.GetAll(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(fetchedExports))
	assert.Equal(t, export.Id, fetchedExports[0].Id)
	assert.Equal(t, export.S3Key, fetchedExports[0].S3Key)
	assert.Equal(t, export.ETag, fetchedExports[0].ETag)
	assert.Equal(t, export.Size, fetchedExports[0].Size)
	assert.Equal(t, export.DeviceId, fetchedExports[0].DeviceId)
	assert.NotEmpty(t, export.CreatedAt)
	assert.NotEmpty(t, export.UpdatedAt)
}
