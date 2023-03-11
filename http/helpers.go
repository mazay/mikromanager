package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strconv"

	"github.com/mazay/mikromanager/db"
	"github.com/mazay/mikromanager/utils"
	"github.com/sirupsen/logrus"
)

var (
	loginTmpl         = path.Join("templates", "login.html")
	baseTmpl          = path.Join("templates", "base.html")
	paginationTmpl    = path.Join("templates", "pagination.html")
	indexTmpl         = path.Join("templates", "index.html")
	deviceDetailsTmpl = path.Join("templates", "device-details.html")
	deviceFormTmpl    = path.Join("templates", "device-form.html")
	credsTmpl         = path.Join("templates", "credentials.html")
	credsFormTmpl     = path.Join("templates", "credentials-form.html")
	erpTmpl           = path.Join("templates", "erp-form.html")
	exportsTmpl       = path.Join("templates", "exports.html")
	exportTmpl        = path.Join("templates", "export.html")
	userFormTmpl      = path.Join("templates", "user-form.html")
	usersTmpl         = path.Join("templates", "users.html")
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

func (dh *dynamicHandler) checkSession(r *http.Request) (*utils.Session, error) {
	var (
		err     error
		session = &utils.Session{}
	)

	c, err := r.Cookie("session_token")
	if err != nil {
		return session, err
	}

	session.Id = c.Value
	err = session.GetById(dh.db)
	if err != nil {
		return session, err
	}

	if session.Expired() {
		return session, fmt.Errorf("Session expired")
	}

	return session, err
}
