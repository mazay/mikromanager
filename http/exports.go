package http

import (
	"net/http"
	"strconv"

	"github.com/mazay/mikromanager/utils"
)

type exportsData struct {
	BackupPath  string
	DeviceId    string
	Count       int
	Exports     []*utils.Export
	Pagination  *Pagination
	CurrentPage int
}

func (dh *dynamicHandler) getExports(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		exports    []*utils.Export
		export     = &utils.Export{}
		data       = &exportsData{BackupPath: dh.backupPath}
		id         = r.URL.Query().Get("id")
		pageId     = r.URL.Query().Get("page_id")
		intPageID  = 1
		pagination = &Pagination{}
	)

	if pageId != "" {
		intPageID, err = strconv.Atoi(pageId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if id != "" {
		data.DeviceId = id
		exports, err = export.GetByDeviceId(dh.db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		exports, err = export.GetAll(dh.db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	data.Count = len(exports)
	chunkedExports := chunkSlice(exports, 10)
	pagination.paginate(*r.URL, intPageID, len(chunkedExports))
	data.Pagination = pagination
	data.CurrentPage = intPageID
	data.Exports = chunkedExports[intPageID-1]
	dh.renderTemplate(w, exportsTmpl, data)
}
