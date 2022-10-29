package http

import (
	"html/template"
	"net/http"
	"path"

	"github.com/mazay/mikromanager/utils"
)

type exportsData struct {
	BackupPath string
	DeviceId   string
	Exports    []*utils.Export
}

func (dh *dynamicHandler) getExports(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		erpTmpl = path.Join("templates", "exports.html")
		export  = &utils.Export{}
		exports = []*utils.Export{}
		data    = &exportsData{BackupPath: dh.backupPath}
	)

	var id = r.URL.Query().Get("id")

	if id != "" {
		data.DeviceId = id
		exports, err = export.GetByDeviceId(dh.db, id)
		dh.logger.Debugf("getting exports for device id %s", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		dh.logger.Debug("getting all exports")
		exports, err = export.GetAll(dh.db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	data.Exports = exports
	// load templates
	tmpl, err := template.New("").Funcs(funcMap).ParseFiles(erpTmpl, baseTmpl)
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
