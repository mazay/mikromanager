package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/mazay/mikromanager/db"
	"go.uber.org/zap"
)

var (
	loginTmpl         = path.Join("templates", "login.html")
	baseTmpl          = path.Join("templates", "base.html")
	paginationTmpl    = path.Join("templates", "pagination.html")
	indexTmpl         = path.Join("templates", "index.html")
	deviceDetailsTmpl = path.Join("templates", "device_details.html")
	deviceFormTmpl    = path.Join("templates", "device_form.html")
	credsTmpl         = path.Join("templates", "credentials.html")
	credsFormTmpl     = path.Join("templates", "credentials_form.html")
	erpTmpl           = path.Join("templates", "erp_form.html")
	exportsTmpl       = path.Join("templates", "exports.html")
	exportTmpl        = path.Join("templates", "export.html")
	userFormTmpl      = path.Join("templates", "user_form.html")
	usersTmpl         = path.Join("templates", "users.html")
)

func handlerWrapper(fn http.HandlerFunc, logger *zap.Logger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/healthz" {
			logger.Info(
				req.URL.String(),
				zap.String("address", req.RemoteAddr),
				zap.String("method", req.Method),
				zap.String("protocol", req.Proto),
				zap.Int64("size", req.ContentLength),
			)
		}

		res.Header().Set("Server", "mikromanager")

		fn(res, req)
	}
}

func (c *HttpConfig) renderTemplate(w http.ResponseWriter, tmplList []string, data any) {
	var err error

	// load templates
	tmpl, err := template.New("").Funcs(funcMap).ParseFiles(tmplList...)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// render the templates
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *HttpConfig) checkSession(r *http.Request) (*db.Session, error) {
	var (
		err     error
		session = &db.Session{}
	)

	cookie, err := r.Cookie("session_token")
	if err != nil {
		return session, err
	}

	session.ID = cookie.Value
	err = session.GetById(c.Db)
	if err != nil {
		return session, err
	}

	if session.Expired() {
		return session, fmt.Errorf("session expired")
	}

	return session, err
}
