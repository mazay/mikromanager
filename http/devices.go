package http

import (
	"net/http"

	"github.com/mazay/mikromanager/utils"
)

type deviceForm struct {
	Id            string
	Address       string
	ApiPort       string
	SshPort       string
	CredentialsId string
	Msg           string
	Credentials   []*utils.Credentials
}

type deviceDetails struct {
	Device     *utils.Device
	Exports    []*utils.Export
	BackupPath string
}

func (dh *dynamicHandler) editDevice(w http.ResponseWriter, r *http.Request) {
	var (
		deviceErr error
		data      = &deviceForm{}
		creds     = &utils.Credentials{}
	)

	credsAll, _ := creds.GetAll(dh.db)
	data.Credentials = credsAll

	if r.Method == "POST" {
		// parse the form
		err := r.ParseForm()
		if err != nil {
			dh.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := r.PostForm.Get("idInput")
		address := r.PostForm.Get("address")
		apiPort := r.PostForm.Get("apiPort")
		sshPort := r.PostForm.Get("sshPort")
		credentialsId := r.PostForm.Get("credentialsId")

		device := &utils.Device{
			Id:            id,
			Address:       address,
			ApiPort:       apiPort,
			SshPort:       sshPort,
			CredentialsId: credentialsId,
		}

		if id == "" {
			// "id" is unset - create new credentials
			deviceErr = device.Create(dh.db)
		} else {
			// "id" is set - update existing credentials
			err := device.GetById(dh.db)
			if err != nil {
				data.Msg = err.Error()
			}
			device.Address = address
			device.ApiPort = apiPort
			device.SshPort = sshPort
			device.CredentialsId = credentialsId
			deviceErr = device.Update(dh.db)
		}

		if deviceErr != nil {
			// return data with errors if validation failed
			data.Id = id
			data.Address = address
			data.ApiPort = apiPort
			data.SshPort = sshPort
			data.CredentialsId = credentialsId
			data.Msg = deviceErr.Error()
		} else {
			http.Redirect(w, r, "/", 302)
			return
		}
	} else {
		// fill in the form if "id" GET parameter set
		id := r.URL.Query().Get("id")
		if id != "" {
			d := &utils.Device{}
			d.Id = id
			err := d.GetById(dh.db)
			if err != nil {
				data.Msg = err.Error()
			} else {
				data.Id = d.Id
				data.Address = d.Address
				data.ApiPort = d.ApiPort
				data.SshPort = d.SshPort
				data.CredentialsId = d.CredentialsId
			}
		}
	}

	dh.renderTemplate(w, deviceFormTmpl, data)
}

func (dh *dynamicHandler) getDevices(w http.ResponseWriter, r *http.Request) {
	var d = &utils.Device{}

	// fetch devices
	deviceList, err := d.GetAll(dh.db)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dh.renderTemplate(w, indexTmpl, deviceList)
}

func (dh *dynamicHandler) getDevice(w http.ResponseWriter, r *http.Request) {
	var device = &utils.Device{}
	var data = &deviceDetails{BackupPath: dh.backupPath}
	var export = &utils.Export{}
	var id = r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Something went wrong, no device ID provided", http.StatusInternalServerError)
		return
	}

	// fetch device from the DB
	device.Id = id
	err := device.GetById(dh.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Device = device

	exports, err := export.GetByDeviceId(dh.db, device.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data.Exports = exports
	dh.renderTemplate(w, deviceDetailsTmpl, data)
}

func (dh *dynamicHandler) deleteDevice(w http.ResponseWriter, r *http.Request) {
	var d = &utils.Device{}

	d.Id = r.URL.Query().Get("id")

	err := d.Delete(dh.db)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", 302)
}
