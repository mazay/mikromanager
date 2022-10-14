package utils

import (
	"encoding/json"
	"log"

	database "github.com/mazay/mikromanager/db"
)

type Device struct {
	Id                   string `json:"_id"`
	Address              string `json:"address"`
	ArchitectureName     string `json:"architecture-name"`
	BadBlocks            int64  `json:"bad-blocks"`
	BoardName            string `json:"board-name"`
	BuildTime            string `json:"build-time"`
	CPU                  string `json:"cpu"`
	CpuCount             int    `json:"cpu-count"`
	CpuFrequency         int    `json:"cpu-frequency"`
	CpuLoad              int    `json:"cpu-load"`
	FactorySoftware      string `json:"factory-software"`
	FreeHddSpace         int64  `json:"free-hdd-space"`
	FreeMemory           int64  `json:"free-memory"`
	Identity             string `json:"identity"`
	Platform             string `json:"platform"`
	TotalHddSpace        int64  `json:"total-hdd-space"`
	TotalMemory          int64  `json:"total-memory"`
	Uptime               string `json:"uptime"`
	Version              string `json:"version"`
	WriteSectSinceReboot int64  `json:"write-sect-since-reboot"`
	WriteSectTotal       int64  `json:"write-sect-total"`
}

func (d *Device) GetAll(db *database.DB) ([]*Device, error) {
	var deviceList []*Device

	docs, err := db.FindAll("devices")
	if err != nil {
		log.Fatal(err)
		return deviceList, err
	}

	for _, doc := range docs {
		dm := &Device{}
		dj, _ := json.Marshal(doc)
		_ = json.Unmarshal(dj, dm)
		deviceList = append(deviceList, dm)
	}

	return deviceList, nil
}
