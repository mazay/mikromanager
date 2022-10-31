package http

import (
	"html/template"
	"net/http"
	"path"
	"strconv"

	"github.com/mazay/mikromanager/db"
	"github.com/sirupsen/logrus"
)

var (
	baseTmpl          = path.Join("templates", "base.html")
	paginationTmpl    = path.Join("templates", "pagination.html")
	indexTmpl         = path.Join("templates", "index.html")
	deviceDetailsTmpl = path.Join("templates", "device-details.html")
	deviceFormTmpl    = path.Join("templates", "device-form.html")
	credsTmpl         = path.Join("templates", "credentials.html")
	credsFormTmpl     = path.Join("templates", "credentials-form.html")
	erpTmpl           = path.Join("templates", "erp-form.html")
	exportsTmpl       = path.Join("templates", "exports.html")
)

type dynamicHandler struct {
	db            *db.DB
	encryptionKey string
	logger        *logrus.Entry
	backupPath    string
}

func handlerWrapper(fn http.HandlerFunc, logger *logrus.Entry) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Infof("%s - \"%s %s %s\" %s", req.RemoteAddr, req.Method,
			req.URL.String(), req.Proto, strconv.FormatInt(req.ContentLength, 10))

		res.Header().Set("Server", "mikromanager")

		fn(res, req)
	}
}

func (dh *dynamicHandler) renderTemplate(w http.ResponseWriter, tmplList []string, data any) {
	var err error

	tmplList = append(tmplList, baseTmpl)

	// load templates
	tmpl, err := template.New("").Funcs(funcMap).ParseFiles(tmplList...)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// render the templates
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
