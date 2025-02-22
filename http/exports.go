package http

import (
	"fmt"
	"net/http"
	"sort"

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
func (c *HttpConfig) getExports(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		device     = &utils.Device{}
		data       = &exportsData{BackupPath: c.BackupPath}
		id         = r.URL.Query().Get("id")
		pagination = &Pagination{}
		templates  = []string{exportsTmpl, paginationTmpl, baseTmpl}
	)

	data.DeviceId = id

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	data.Devices, err = device.GetAll(c.Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageId, perPage, err := getPagionationParams(r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exports, err := c.S3.GetExports(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// sort exports by last modified date in descending order
	sort.Slice(exports, func(i, j int) bool {
		return exports[i].LastModified.After(*exports[j].LastModified)
	})

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

	c.renderTemplate(w, templates, data)
}

// getExport responds to GET /exports?id=<id> and displays the export by <id>
func (c *HttpConfig) getExport(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		device    = &utils.Device{}
		data      = &exportData{BackupPath: c.BackupPath}
		id        = r.URL.Query().Get("id")
		templates = []string{exportTmpl, baseTmpl}
	)

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if id == "" {
		http.Error(w, "Export not found", http.StatusNotFound)
		return
	}

	export, err := c.S3.GetExportAttributes(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Export = export

	device.Id = export.GetDeviceId()
	err = device.GetById(c.Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Device = device

	exportBody, err := export.GetBody(c.S3)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.ExportData = string(exportBody)

	c.renderTemplate(w, templates, data)
}

// downloadExport responds to GET /exports/download?id=<id> and downloads the export by <id>
func (c *HttpConfig) downloadExport(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		device = &utils.Device{}
		id     = r.URL.Query().Get("id")
	)

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if id == "" {
		http.Error(w, "Export not found", http.StatusNotFound)
		return
	}

	export, err := c.S3.GetExportAttributes(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	device.Id = export.GetDeviceId()
	err = device.GetById(c.Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exportBody, err := export.GetBody(c.S3)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%s %s.rsc", device.Identity, export.LastModified.Format("2006-01-02 15:04:05"))

	// stream the export file
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write(exportBody)
	if err != nil {
		c.Logger.Error(err.Error())
	}
}
