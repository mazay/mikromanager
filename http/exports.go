package http

import (
	"fmt"
	"net/http"

	"github.com/mazay/mikromanager/internal"
	"github.com/mazay/mikromanager/utils"
)

type exportsData struct {
	BackupPath  string
	DeviceId    string
	Count       int
	Exports     []*internal.Export
	Devices     []*utils.Device
	Pagination  *Pagination
	CurrentPage int
}

type exportData struct {
	BackupPath string
	Export     *internal.Export
	Device     *utils.Device
	ExportData string
}

// getExports handles the GET request for /exports and displays a paginated list of exports
// for the specified device ID. It retrieves all devices and exports, applies pagination based
// on query parameters, and renders the exports template with the gathered data.
func (dh *dynamicHandler) getExports(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		device     = &utils.Device{}
		data       = &exportsData{BackupPath: dh.backupPath}
		id         = r.URL.Query().Get("id")
		pagination = &Pagination{}
		templates  = []string{exportsTmpl, paginationTmpl, baseTmpl}
	)

	data.DeviceId = id

	_, err = dh.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
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
	exports, err := dh.s3.GetExports(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

// getExport responds to GET /exports?id=<id> and displays the export by <id>
func (dh *dynamicHandler) getExport(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		device    = &utils.Device{}
		data      = &exportData{BackupPath: dh.backupPath}
		id        = r.URL.Query().Get("id")
		templates = []string{exportTmpl, baseTmpl}
	)

	_, err = dh.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if id == "" {
		http.Error(w, "Export not found", http.StatusNotFound)
		return
	}

	export, err := dh.s3.GetExport(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Export = export

	device.Id = export.GetDeviceId()
	err = device.GetById(dh.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Device = device

	exportBody, err := export.GetBody(dh.s3)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.ExportData = string(exportBody)

	dh.renderTemplate(w, templates, data)
}

// downloadExport responds to GET /exports/download?id=<id> and downloads the export by <id>
func (dh *dynamicHandler) downloadExport(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		device = &utils.Device{}
		id     = r.URL.Query().Get("id")
	)

	_, err = dh.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if id == "" {
		http.Error(w, "Export not found", http.StatusNotFound)
		return
	}

	export, err := dh.s3.GetExport(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	device.Id = export.GetDeviceId()
	err = device.GetById(dh.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exportBody, err := export.GetBody(dh.s3)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%s %s.rsc", device.Identity, export.LastModified.Format("2006-01-02 15:04:05"))

	// stream the export file
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "text/plain")
	w.Write(exportBody)
}
