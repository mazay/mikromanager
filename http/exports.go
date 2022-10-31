package http

import (
	"net/http"

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
		pagination = &Pagination{}
	)

	err, pageId, perPage := getPagionationParams(r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	chunkedExports := chunkSliceOfObjects(exports, perPage)
	pagination.paginate(*r.URL, pageId, len(chunkedExports))

	if pageId-1 >= len(chunkedExports) {
		pageId = len(chunkedExports)
	}

	data.Pagination = pagination
	data.CurrentPage = pageId
	data.Exports = chunkedExports[pageId-1]
	dh.renderTemplate(w, []string{exportsTmpl, paginationTmpl}, data)
}
