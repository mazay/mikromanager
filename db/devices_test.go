package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testDevices = []*Device{
		{
			Address:              "10.0.0.1",
			ApiPort:              "8728",
			ArchitectureName:     "arm64",
			BadBlocks:            0,
			BoardName:            "RB5009UPr+S+",
			BuildTime:            "2018-01-01 00:00:00",
			CPU:                  "ARM64",
			CpuCount:             4,
			CpuFrequency:         1400,
			CpuLoad:              10,
			FactorySoftware:      "7.15.2",
			FreeHddSpace:         1000,
			FreeMemory:           256,
			Identity:             "RB5009UPr",
			Platform:             "MikroTik",
			PolledAt:             time.Now().UTC(),
			PollingSucceeded:     1,
			SshPort:              "22",
			TotalHddSpace:        1000,
			TotalMemory:          1000,
			Uptime:               "1w2d3h4m5s",
			Version:              "7.17.2",
			WriteSectSinceReboot: 10,
			WriteSectTotal:       20,
			Model:                "RB5009UPr+S+",
			SerialNumber:         "1234567890",
			FirmwareType:         "70x0",
			FactoryFirmware:      "7.15.2",
			CurrentFirmware:      "7.17.2",
			UpgradeFirmware:      "7.17.2",
		},
		{
			Address: "10.0.0.2",
		},
	}
)

func TestDevicesGetAllPlain(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	group, err := createTestDeviceGroup(db)
	if err != nil {
		t.Fatal(err)
	}

	for _, dev := range testDevices {
		dev.Groups = []*DeviceGroup{group}
		err = dev.Create(db)
		if err != nil {
			t.Fatal(err)
		}
	}

	d := &Device{}
	fetchedDevs, err := d.GetAllPlain(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, len(fetchedDevs))
	assert.Len(t, fetchedDevs[0].Groups, 0)
	assert.Equal(t, testDevices[0].Address, fetchedDevs[0].Address)
	assert.Equal(t, testDevices[0].ApiPort, fetchedDevs[0].ApiPort)
	assert.Equal(t, testDevices[0].ArchitectureName, fetchedDevs[0].ArchitectureName)
	assert.Equal(t, testDevices[0].BadBlocks, fetchedDevs[0].BadBlocks)
	assert.Equal(t, testDevices[0].BoardName, fetchedDevs[0].BoardName)
	assert.Equal(t, testDevices[0].BuildTime, fetchedDevs[0].BuildTime)
	assert.Equal(t, testDevices[0].CPU, fetchedDevs[0].CPU)
	assert.Equal(t, testDevices[0].CpuCount, fetchedDevs[0].CpuCount)
	assert.Equal(t, testDevices[0].CpuFrequency, fetchedDevs[0].CpuFrequency)
	assert.Equal(t, testDevices[0].CpuLoad, fetchedDevs[0].CpuLoad)
	assert.Equal(t, testDevices[0].FactorySoftware, fetchedDevs[0].FactorySoftware)
	assert.Equal(t, testDevices[0].FreeHddSpace, fetchedDevs[0].FreeHddSpace)
	assert.Equal(t, testDevices[0].FreeMemory, fetchedDevs[0].FreeMemory)
	assert.Equal(t, testDevices[0].Identity, fetchedDevs[0].Identity)
	assert.Equal(t, testDevices[0].Platform, fetchedDevs[0].Platform)
	assert.Equal(t, testDevices[0].PolledAt, fetchedDevs[0].PolledAt)
	assert.Equal(t, testDevices[0].PollingSucceeded, fetchedDevs[0].PollingSucceeded)
	assert.Equal(t, testDevices[0].SshPort, fetchedDevs[0].SshPort)
	assert.Equal(t, testDevices[0].TotalHddSpace, fetchedDevs[0].TotalHddSpace)
	assert.Equal(t, testDevices[0].TotalMemory, fetchedDevs[0].TotalMemory)
	assert.Equal(t, testDevices[0].Uptime, fetchedDevs[0].Uptime)
	assert.Equal(t, testDevices[0].Version, fetchedDevs[0].Version)
	assert.Equal(t, testDevices[0].WriteSectSinceReboot, fetchedDevs[0].WriteSectSinceReboot)
	assert.Equal(t, testDevices[0].WriteSectTotal, fetchedDevs[0].WriteSectTotal)
	assert.Equal(t, testDevices[0].Model, fetchedDevs[0].Model)
	assert.Equal(t, testDevices[0].SerialNumber, fetchedDevs[0].SerialNumber)
	assert.Equal(t, testDevices[0].FirmwareType, fetchedDevs[0].FirmwareType)
	assert.Equal(t, testDevices[0].FactoryFirmware, fetchedDevs[0].FactoryFirmware)
	assert.Equal(t, testDevices[0].CurrentFirmware, fetchedDevs[0].CurrentFirmware)
	assert.Equal(t, testDevices[0].UpgradeFirmware, fetchedDevs[0].UpgradeFirmware)
	assert.Equal(t, testDevices[1].Address, fetchedDevs[1].Address)

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDevicesGetCredentials(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	creds := []*Credentials{
		{
			Alias:             "Default",
			Username:          "test-username",
			EncryptedPassword: "test-password",
		},
		{
			Alias:             "test-alias",
			Username:          "test-username2",
			EncryptedPassword: "test-password2",
		},
	}

	for _, c := range creds {
		err = c.Create(db)
		if err != nil {
			t.Fatal(err)
		}
	}

	devs := []*Device{
		// device uses non-default credentials
		{
			Address:     "10.0.0.1",
			Credentials: creds[1],
		},
		// device uses default credentials
		{
			Address: "10.0.0.2",
		},
	}

	for _, d := range devs {
		err = d.Create(db)
		if err != nil {
			t.Fatal(err)
		}

		fetchedCreds, err := d.GetCredentials(db)
		if err != nil {
			t.Fatal(err)
		}

		if d.Credentials == nil {
			assert.Equal(t, "Default", fetchedCreds.Alias)
		} else {
			assert.Equal(t, "test-alias", fetchedCreds.Alias)
		}
	}
}

func TestDevicesCreate(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	dev := &Device{
		Address: "10.0.0.1",
	}

	err = dev.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, dev.Id)
	assert.NotEmpty(t, dev.CreatedAt)
	assert.NotEmpty(t, dev.UpdatedAt)

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDevicesUpdate(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	dev := &Device{
		Address: "10.0.0.1",
	}

	err = dev.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	dev.Address = "10.0.0.2"
	err = dev.Update(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "10.0.0.2", dev.Address)

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDevicesGetById(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	dev := &Device{
		Address: "10.0.0.1",
	}

	err = dev.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	fetchedDev := &Device{}
	fetchedDev.Id = dev.Id
	err = fetchedDev.GetById(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "10.0.0.1", fetchedDev.Address)
	assert.NotEmpty(t, fetchedDev.Id)
	assert.NotEmpty(t, fetchedDev.CreatedAt)
	assert.NotEmpty(t, fetchedDev.UpdatedAt)

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDevicesDelete(t *testing.T) {
	db, err := openTestDb(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	dev := &Device{
		Address: "10.0.0.1",
	}

	err = dev.Create(db)
	if err != nil {
		t.Fatal(err)
	}

	err = dev.Delete(db)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
}
