package internal

import (
	"encoding/json"

	"github.com/mazay/mikromanager/db"
	"go.uber.org/zap"
)

// CheckForUpdates fetches the latest routerboard update information for the given device.
//
// The function will query the device using the Mikrotik API to check for and fetch the latest update
// information. It will then marshal the response into a JSON byte slice and unmarshal it into the
// device object. The function will return an error if any of the queries or marshaling/unmarshaling
// operations fail.
func (api *Api) CheckForUpdates(device *db.Device) error {
	_, err := api.Run("/system/package/update/check-for-updates ?once")
	if err != nil {
		return err
	}
	resource, err := api.Run("/system/package/update/getall")
	if err != nil {
		return err
	}
	api.Logger.Debug("fetched rouerboard update data", zap.String("address", api.Address))
	inrec, _ := json.Marshal(resource[0].Map)
	return json.Unmarshal(inrec, &device)
}

// UpdateDevice will update the specified Mikrotik device using the Mikrotik API.
//
// The function will first retrieve the credentials associated with the device from the database
// and decrypt the password using the provided encryption key. It will then create a new Mikrotik
// API client with the decrypted credentials and attempt to update the device using the API.
// If any of the operations fail, the function will log an error message using the provided logger.
func UpdateDevice(device *db.Device, database *db.DB, encryptionKey string, logger *zap.Logger) {
	var err error
	credentials, err := device.GetCredentials(database)
	if err != nil {
		logger.Error(err.Error())
	}

	decryptedPw, err := db.DecryptString(credentials.EncryptedPassword, encryptionKey)
	if err != nil {
		logger.Error(err.Error())
	}

	client := &Api{
		Address:  device.Address,
		Port:     device.ApiPort,
		Username: credentials.Username,
		Password: decryptedPw,
		Async:    false,
		UseTLS:   false,
	}
	_, err = client.Run("/system/package/update/install")
	if err != nil {
		logger.Error(err.Error())
	}
}
