package http

import (
	"html/template"
	"net/http"
	"path"

	"github.com/mazay/mikromanager/utils"
)

func (dh *dynamicHandler) getDevices(w http.ResponseWriter, r *http.Request) {
	var indexTmpl = path.Join("templates", "index.html")
	var d = &utils.Device{}

	// fetch devices
	deviceList, err := d.GetAll(dh.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// load templates
	tmpl, err := template.New("").ParseFiles(indexTmpl, baseTmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// render the templates
	if err := tmpl.ExecuteTemplate(w, "base", deviceList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
