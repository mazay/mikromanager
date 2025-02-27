package internal

import (
	"encoding/json"
	"fmt"

	"github.com/go-routeros/routeros/v3/proto"
	"github.com/mazay/mikromanager/db"
)

type CpuResourceItem struct {
	Id   string `json:".id"`
	Core string `json:"cpu"`
	Load string `json:"load"`
	Irq  string `json:"irq"`
	Disk string `json:"disk"`
}

type CpuResource struct {
	Load string
	Irq  string
	Disk string
}

// Parse takes a slice of routeros sentences and returns a slice of CpuResourceItem objects.
// It extracts the CPU resource information from each sentence and appends it to the items slice.
// Returns an error if JSON unmarshalling fails for any of the sentences.
func (c *CpuResource) Parse(outputs []*proto.Sentence) ([]*CpuResourceItem, error) {
	var (
		err   error
		items []*CpuResourceItem
	)
	for _, i := range outputs {
		item := CpuResourceItem{}
		inrec, _ := json.Marshal(i.Map)
		err = json.Unmarshal(inrec, &item)
		if err != nil {
			return items, err
		}
		items = append(items, &item)
	}
	return items, nil
}

// SentencesToCpuResources converts a slice of routeros sentences into a map of CpuResource objects.
// Each sentence is parsed into a CpuResourceItem, and the function builds a map using the CPU core
// identifier as the key. Returns an error if the input slice is empty, if parsing fails, or if any
// CpuResourceItem is missing a core identifier.
func SentencesToCpuResources(outputs []*proto.Sentence) (map[string]*CpuResource, error) {
	if len(outputs) == 0 {
		return nil, fmt.Errorf("empty outputs")
	}

	result := make(map[string]*CpuResource)
	hi := &CpuResource{}
	items, err := hi.Parse(outputs)
	if err != nil {
		return nil, err
	}

	for _, i := range items {
		if i.Core == "" {
			return nil, fmt.Errorf("missing core identifier")
		}
		result[i.Core] = &CpuResource{
			Load: i.Load,
			Irq:  i.Irq,
			Disk: i.Disk,
		}
	}

	return result, nil
}

// GetCpuResources retrieves the CPU resources from a Mikrotik device and saves them to the database.
// It executes the "/system/resource/cpu/getall" command via the Mikrotik API, parses the response into a map of
// CpuResource objects, and returns the map. If any step fails, an error is returned.
func GetCpuResources(device *db.Device, database *db.DB, encryptionKey string) (map[string]*CpuResource, error) {
	credentials, err := device.GetCredentials(database)
	if err != nil {
		return nil, err
	}

	encryptedPassword, err := db.DecryptString(credentials.EncryptedPassword, encryptionKey)
	if err != nil {
		return nil, err
	}

	api := &Api{
		Address:  device.Address,
		Username: credentials.Username,
		Password: encryptedPassword,
		Async:    true,
	}
	resource, err := api.Run("/system/resource/cpu/getall")
	if err != nil {
		return nil, err
	}

	return SentencesToCpuResources(resource)
}
