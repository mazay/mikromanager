package http

import (
	"net/http"

	"github.com/mazay/mikromanager/db"
	"github.com/mazay/mikromanager/internal"
)

type deviceForm struct {
	Id            string
	Address       string
	ApiPort       string
	SshPort       string
	CredentialsId string
	Msg           string
	Credentials   []*db.Credentials
}

type deviceDetails struct {
	Device  *db.Device
	Exports []*db.Export
}

type devicesData struct {
	Count       int
	Devices     []*db.Device
	Pagination  *Pagination
	CurrentPage int
}

func (df *deviceForm) formFillIn(device *db.Device) {
	df.Id = device.Id
	df.Address = device.Address
	df.ApiPort = device.ApiPort
	df.SshPort = device.SshPort
	df.CredentialsId = device.CredentialsID
}

func (c *HttpConfig) editDevice(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		deviceErr error
		data      = &deviceForm{}
		creds     = &db.Credentials{}
		templates = []string{deviceFormTmpl, baseTmpl}
	)

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	credsAll, err := creds.GetAll(c.Db)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Credentials = credsAll

	if r.Method == "POST" {
		// parse the form
		err = r.ParseForm()
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := r.PostForm.Get("idInput")
		address := r.PostForm.Get("address")
		apiPort := r.PostForm.Get("apiPort")
		sshPort := r.PostForm.Get("sshPort")
		credentialsId := r.PostForm.Get("credentialsId")

		device := &db.Device{
			Address: address,
			ApiPort: apiPort,
			SshPort: sshPort,
		}

		device.Id = id
		device.CredentialsID = credentialsId

		if id == "" {
			// "id" is unset - create new credentials
			deviceErr = device.Create(c.Db)
		} else {
			// "id" is set - update existing credentials
			err := device.GetById(c.Db)
			if err != nil {
				data.Msg = err.Error()
			}
			device.Address = address
			device.ApiPort = apiPort
			device.SshPort = sshPort
			device.CredentialsID = credentialsId
			deviceErr = device.Update(c.Db)
		}

		if deviceErr != nil {
			// return data with errors if validation failed
			data.Id = id
			data.formFillIn(device)
			data.Msg = deviceErr.Error()
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		// fill in the form if "id" GET parameter set
		id := r.URL.Query().Get("id")
		if id != "" {
			d := &db.Device{}
			d.Id = id
			err = d.GetById(c.Db)
			if err != nil {
				data.Msg = err.Error()
			} else {
				data.formFillIn(d)
			}
		}
	}

	c.renderTemplate(w, templates, data)
}

func (c *HttpConfig) getDevices(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		d          = &db.Device{}
		data       = &devicesData{}
		pagination = &Pagination{}
		templates  = []string{indexTmpl, paginationTmpl, baseTmpl, updateModalTmpl}
	)

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	pageId, perPage, err := getPagionationParams(r.URL)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fetch devices
	deviceList, err := d.GetAllPreload(c.Db)
	if err != nil {
		c.Logger.Error(err.Error())
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

	c.renderTemplate(w, templates, data)
}

func (c *HttpConfig) getDevice(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		device    = &db.Device{}
		export    = &db.Export{}
		data      = &deviceDetails{}
		id        = r.URL.Query().Get("id")
		templates = []string{deviceDetailsTmpl, baseTmpl, updateModalTmpl}
	)

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if id == "" {
		http.Error(w, "Something went wrong, no device ID provided", http.StatusInternalServerError)
		return
	}

	// fetch device from the DB
	device.Id = id
	err = device.GetById(c.Db)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Device = device

	exports, err := export.GetByDeviceId(c.Db, id)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Exports = exports
	c.renderTemplate(w, templates, data)
}

func (c *HttpConfig) deleteDevice(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		d   = &db.Device{}
		id  = r.URL.Query().Get("id")
	)

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if id == "" {
		http.Error(w, "Something went wrong, no device ID provided", http.StatusInternalServerError)
		return
	}

	d.Id = id

	// delete device
	err = d.Delete(c.Db)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// delete exports
	exports, err := c.S3.GetExports(d.Id)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = c.S3.DeleteExports(exports)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (c *HttpConfig) updateDevice(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		d   = &db.Device{}
		id  = r.URL.Query().Get("id")
	)

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if id == "" {
		http.Error(w, "Something went wrong, no device ID provided", http.StatusInternalServerError)
		return
	}

	d.Id = id
	err = d.GetById(c.Db)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// trigger an update without blocking
	go internal.UpdateDevice(d, c.Db, c.EncryptionKey, c.Logger)
	w.WriteHeader(http.StatusOK)
}
