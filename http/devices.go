package http

import (
	"html/template"
	"io"
	"net/http"
	"path"

	"github.com/mazay/mikromanager/utils"
)

func (dh *dynamicHandler) getDevices(w http.ResponseWriter, r *http.Request) {
	var indexTmpl = path.Join("templates", "index.html")
	var d = &utils.Device{}

	devices, err := dh.db.FindAll("devices")
	if err != nil {
		io.WriteString(w, err.Error())
	}

	deviceList := d.FromListOfMaps(devices)

	tmpl, err := template.New("").ParseFiles(indexTmpl, baseTmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", deviceList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
