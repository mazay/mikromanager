package http

import (
	"html/template"
	"net/http"
	"path"

	"github.com/mazay/mikromanager/utils"
)

type deviceForm struct {
	Id            string
	Address       string
	Port          string
	CredentialsId string
	Msg           string
	Credentials   []*utils.Credentials
}

func (dh *dynamicHandler) editDevice(w http.ResponseWriter, r *http.Request) {
	var (
		deviceErr  error
		deviceTmpl = path.Join("templates", "device-form.html")
		data       = &deviceForm{}
		creds      = &utils.Credentials{}
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
		port := r.PostForm.Get("port")
		credentialsId := r.PostForm.Get("credentialsId")

		device := &utils.Device{
			Id:            id,
			Address:       address,
			Port:          port,
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
			device.Port = port
			device.CredentialsId = credentialsId
			deviceErr = device.Update(dh.db)
		}

		if deviceErr != nil {
			// return data with errors if validation failed
			data.Id = id
			data.Address = address
			data.Port = port
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
				data.Port = d.Port
				data.CredentialsId = d.CredentialsId
			}
		}
	}

	// load templates
	tmpl, err := template.New("").ParseFiles(deviceTmpl, baseTmpl)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// render the templates
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (dh *dynamicHandler) getDevices(w http.ResponseWriter, r *http.Request) {
	var indexTmpl = path.Join("templates", "index.html")
	var d = &utils.Device{}

	// fetch devices
	deviceList, err := d.GetAll(dh.db)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// load templates
	tmpl, err := template.New("").Funcs(funcMap).ParseFiles(indexTmpl, baseTmpl)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// render the templates
	if err := tmpl.ExecuteTemplate(w, "base", deviceList); err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (dh *dynamicHandler) getDevice(w http.ResponseWriter, r *http.Request) {
	var deviceTmpl = path.Join("templates", "device-details.html")
	var device = &utils.Device{}
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

	// load templates
	tmpl, err := template.New("").Funcs(funcMap).ParseFiles(deviceTmpl, baseTmpl)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// render the templates
	if err := tmpl.ExecuteTemplate(w, "base", device); err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
