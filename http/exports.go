package http

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mazay/mikromanager/utils"
)

type exportsData struct {
	BackupPath  string
	DeviceId    string
	Count       int
	Exports     []*utils.Export
	Devices     []*utils.Device
	Pagination  *Pagination
	CurrentPage int
}

type exportData struct {
	BackupPath string
	Export     *utils.Export
	Device     *utils.Device
	ExportData string
}

func (dh *dynamicHandler) getExports(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		exports    []*utils.Export
		export     = &utils.Export{}
		device     = &utils.Device{}
		data       = &exportsData{BackupPath: dh.backupPath}
		id         = r.URL.Query().Get("id")
		pagination = &Pagination{}
		templates  = []string{exportsTmpl, paginationTmpl, baseTmpl}
	)

	_, err = dh.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	data.Devices, err = device.GetAll(dh.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err, pageId, perPage := getPagionationParams(r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dh.db.Sort("created", -1)
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
	if data.Count > 0 {
		chunkedExports := chunkSliceOfObjects(exports, perPage)
		pagination.paginate(*r.URL, pageId, len(chunkedExports))

		if pageId-1 >= len(chunkedExports) {
			pageId = len(chunkedExports)
		}

		data.Pagination = pagination
		data.CurrentPage = pageId
		data.Exports = chunkedExports[pageId-1]
	}

	dh.renderTemplate(w, templates, data)
}

func (dh *dynamicHandler) getExport(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		export    = &utils.Export{}
		device    = &utils.Device{}
		data      = &exportData{BackupPath: dh.backupPath}
		id        = r.URL.Query().Get("id")
		templates = []string{exportTmpl, baseTmpl}
	)

	_, err = dh.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	if id == "" {
		http.Error(w, "Export not found", http.StatusNotFound)
		return
	}

	export.Id = id
	err = export.GetById(dh.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Export = export

	device.Id = export.DeviceId
	err = device.GetById(dh.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Device = device

	exportBody, err := os.ReadFile(data.Export.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.ExportData = string(exportBody)

	dh.renderTemplate(w, templates, data)
}

func (dh *dynamicHandler) downloadExport(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		export = &utils.Export{}
		device = &utils.Device{}
		id     = r.URL.Query().Get("id")
	)

	_, err = dh.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	if id == "" {
		http.Error(w, "Export not found", http.StatusNotFound)
		return
	}

	export.Id = id
	err = export.GetById(dh.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	device.Id = export.DeviceId
	err = device.GetById(dh.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// generate human friendly filename
	filename := fmt.Sprintf("%s %s.rsc", device.Identity, export.Created.Format("2006-01-02 15:04:05"))

	// get file info
	fileInfo, err := os.Stat(export.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileReader, err := os.Open(export.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprint(fileInfo.Size()))

	// stream the body to the client
	_, err = io.Copy(w, fileReader)
	if err != nil {
		dh.logger.Error(err)
	}
}
