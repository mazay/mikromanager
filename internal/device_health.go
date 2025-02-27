package internal

import (
	"encoding/json"
	"fmt"

	"github.com/go-routeros/routeros/v3/proto"
	"github.com/mazay/mikromanager/db"
	"go.uber.org/zap"
)

type HealthItem struct {
	Id    string `json:".id"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

// Parse takes a slice of routeros sentences and returns a slice of HealthItem objects.
// It extracts the health information from the sentences and returns an error if the
// slice of sentences is empty or if any of the expected health items are missing.
func (hi *HealthItem) Parse(outputs []*proto.Sentence) ([]*HealthItem, error) {
	var (
		err   error
		items []*HealthItem
	)

	for _, i := range outputs {
		item := HealthItem{}
		inrec, _ := json.Marshal(i.Map)
		err = json.Unmarshal(inrec, &item)
		if err != nil {
			return items, err
		}
		items = append(items, &item)
	}
	return items, nil
}

// SentencesToHealth takes a slice of routeros sentences and returns a db.Health object. It
// parses the sentences into a slice of HealthItem objects and then extracts the health
// information from those objects. It returns an error if the slice of sentences is empty or
// if any of the expected health items are missing.
func SentencesToHealth(outputs []*proto.Sentence) (*db.Health, error) {
	health := &db.Health{}
	if len(outputs) == 0 {
		return nil, fmt.Errorf("empty output")
	}

	hi := &HealthItem{}
	items, err := hi.Parse(outputs)
	if err != nil {
		return nil, err
	}

	for _, i := range items {
		switch i.Name {
		case "voltage":
			if i.Value == "" {
				return nil, fmt.Errorf("voltage is empty")
			}
			health.Voltage, err = parseFloat32(i.Value)
			if err != nil {
				return nil, err
			}
		case "cpu-temperature":
			if i.Value == "" {
				return nil, fmt.Errorf("cpu-temperature is empty")
			}
			health.CpuTemp, err = parseFloat32(i.Value)
			if err != nil {
				return nil, err
			}
		case "board-temperature1":
			if i.Value == "" {
				return nil, fmt.Errorf("board-temperature1 is empty")
			}
			health.BoardTemp1, err = parseFloat32(i.Value)
			if err != nil {
				return nil, err
			}
		case "board-temperature2":
			if i.Value == "" {
				return nil, fmt.Errorf("board-temperature2 is empty")
			}
			health.BoardTemp2, err = parseFloat32(i.Value)
			if err != nil {
				return nil, err
			}
		case "sfp-temperature":
			if i.Value == "" {
				return nil, fmt.Errorf("sfp-temperature is empty")
			}
			health.SfpTemp, err = parseFloat32(i.Value)
			if err != nil {
				return nil, err
			}
		case "fan-state":
			if i.Value == "" {
				return nil, fmt.Errorf("fan-state is empty")
			}
			health.FanState = i.Value
		case "fan-speed":
			if i.Value == "" {
				return nil, fmt.Errorf("fan-speed is empty")
			}
			health.FanSpeed, err = parseInt(i.Value)
			if err != nil {
				return nil, err
			}
		case "psu1-voltage":
			if i.Value == "" {
				return nil, fmt.Errorf("psu1-voltage is empty")
			}
			health.Psu1Voltage, err = parseFloat32(i.Value)
			if err != nil {
				return nil, err
			}
		case "psu2-voltage":
			if i.Value == "" {
				return nil, fmt.Errorf("psu2-voltage is empty")
			}
			health.Psu2Voltage, err = parseFloat32(i.Value)
			if err != nil {
				return nil, err
			}
		case "poe-out-consumption":
			if i.Value == "" {
				return nil, fmt.Errorf("poe-out-consumption is empty")
			}
			consumption, err := parseFloat32(i.Value)
			if err != nil {
				return nil, err
			}
			health.PoeOutConsumption = consumption
		case "jack-voltage":
			if i.Value == "" {
				return nil, fmt.Errorf("jack-voltage is empty")
			}
			health.JackVoltage, err = parseFloat32(i.Value)
			if err != nil {
				return nil, err
			}
		case "2pin-voltage":
			if i.Value == "" {
				return nil, fmt.Errorf("2pin-voltage is empty")
			}
			health.TwoPinVoltage, err = parseFloat32(i.Value)
			if err != nil {
				return nil, err
			}
		case "poe-in-voltage":
			if i.Value == "" {
				return nil, fmt.Errorf("poe-in-voltage is empty")
			}
			health.PoeInVoltage, err = parseFloat32(i.Value)
			if err != nil {
				return nil, err
			}
		}
	}
	return health, nil
}

// GetDeviceHealth retrieves the health information from a MikroTik device
// and saves it to the database. It executes the "/system/health/getall"
// command via the Mikrotik API, parses the response into a Health object,
// and then saves the health data associated with the given device to the
// specified database. If any step fails, an error is logged and the
// process is aborted.
func (api *Api) GetDeviceHealth(device *db.Device, database *db.DB) {
	resource, err := api.Run("/system/health/getall")
	if err != nil {
		api.Logger.Error("failed to fetch health data", zap.Error(err))
		return
	}
	api.Logger.Debug("fetched rouerboard health data", zap.String("address", api.Address))
	health, err := SentencesToHealth(resource)
	if err != nil {
		api.Logger.Error("failed to parse health data", zap.Error(err))
		return
	}
	health.DeviceId = device.Id
	err = health.Save(database)
	if err != nil {
		api.Logger.Error("failed to save health data", zap.Error(err))
	}
}
