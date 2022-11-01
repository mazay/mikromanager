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

type devicesData struct {
	Count       int
	Devices     []*utils.Device
	Pagination  *Pagination
	CurrentPage int
}

func (df *deviceForm) formFillIn(device *utils.Device) {
	df.Id = device.Id
	df.Address = device.Address
	df.ApiPort = device.ApiPort
	df.SshPort = device.SshPort
	df.CredentialsId = device.CredentialsId
}

func (dh *dynamicHandler) editDevice(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		deviceErr error
		data      = &deviceForm{}
		creds     = &utils.Credentials{}
		templates = []string{deviceFormTmpl, baseTmpl}
	)

	_, err = dh.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	credsAll, err := creds.GetAll(dh.db)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Credentials = credsAll

	if r.Method == "POST" {
		// parse the form
		err = r.ParseForm()
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
			data.formFillIn(device)
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
			err = d.GetById(dh.db)
			if err != nil {
				data.Msg = err.Error()
			} else {
				data.formFillIn(d)
			}
		}
	}

	dh.renderTemplate(w, templates, data)
}

func (dh *dynamicHandler) getDevices(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		d          = &utils.Device{}
		data       = &devicesData{}
		pagination = &Pagination{}
		templates  = []string{indexTmpl, paginationTmpl, baseTmpl}
	)

	_, err = dh.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	err, pageId, perPage := getPagionationParams(r.URL)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fetch devices
	deviceList, err := d.GetAll(dh.db)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data.Count = len(deviceList)
	if data.Count > 0 {
		chunkedDevices := chunkSliceOfObjects(deviceList, perPage)
		pagination.paginate(*r.URL, pageId, len(chunkedDevices))

		if pageId-1 >= len(chunkedDevices) {
			pageId = len(chunkedDevices)
		}
		data.Pagination = pagination
		data.CurrentPage = pageId
		data.Devices = chunkedDevices[pageId-1]
	}

	dh.renderTemplate(w, templates, data)
}

func (dh *dynamicHandler) getDevice(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		device    = &utils.Device{}
		data      = &deviceDetails{BackupPath: dh.backupPath}
		export    = &utils.Export{}
		id        = r.URL.Query().Get("id")
		templates = []string{deviceDetailsTmpl, baseTmpl}
	)

	_, err = dh.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	if id == "" {
		http.Error(w, "Something went wrong, no device ID provided", http.StatusInternalServerError)
		return
	}

	// fetch device from the DB
	device.Id = id
	err = device.GetById(dh.db)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Device = device

	exports, err := export.GetByDeviceId(dh.db, device.Id)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data.Exports = exports
	dh.renderTemplate(w, templates, data)
}

func (dh *dynamicHandler) deleteDevice(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		d   = &utils.Device{}
		id  = r.URL.Query().Get("id")
	)

	_, err = dh.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	if id == "" {
		http.Error(w, "Something went wrong, no device ID provided", http.StatusInternalServerError)
		return
	}

	d.Id = id

	err = d.Delete(dh.db)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", 302)
}
