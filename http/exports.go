package http

import (
	"fmt"
	"net/http"

	"github.com/mazay/mikromanager/db"
)

type exportsData struct {
	DeviceId    string
	Exports     []*db.Export
	Pagination  *Pagination
	CurrentPage int
}

type exportData struct {
	Export     *db.Export
	ExportData string
}

// getExports handles the GET request for /exports and displays a paginated list of exports
// for the specified device ID. It retrieves all devices and exports, applies pagination based
// on query parameters, and renders the exports template with the gathered data.
func (c *HttpConfig) getExports(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		exports    []*db.Export
		export     = &db.Export{}
		data       = &exportsData{}
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

	pageId, perPage, err := getPagionationParams(r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if id != "" {
		exports, err = export.GetByDeviceId(c.Db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		exports, err = export.GetAll(c.Db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if len(exports) > 0 {
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
		export    = &db.Export{}
		data      = &exportData{}
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

	export.Id = id
	err = export.GetById(c.Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Export = export

	exportBody, err := c.S3.GetFile(export.S3Key, *export.Size)
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
		export = &db.Export{}
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

	export.Id = id
	err = export.GetById(c.Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exportBody, err := c.S3.GetFile(export.S3Key, *export.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%s %s.rsc", export.Device.Identity, export.LastModified.Format("2006-01-02 15:04:05"))

	// stream the export file
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write(exportBody)
	if err != nil {
		c.Logger.Error(err.Error())
	}
}
