package http

import (
	"net/http"

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
		exports []*utils.Export
		export  = &utils.Export{}
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
	dh.renderTemplate(w, exportsTmpl, data)
}
