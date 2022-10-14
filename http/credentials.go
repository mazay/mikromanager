package http

import (
	"html/template"
	"net/http"
	"path"

	"github.com/mazay/mikromanager/utils"
)

func (dh *dynamicHandler) getCredentials(w http.ResponseWriter, r *http.Request) {
	var indexTmpl = path.Join("templates", "credentials.html")
	var c = &utils.Credentials{}

	// fetch devices
	credList, err := c.GetAll(dh.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// load templates
	tmpl, err := template.New("").ParseFiles(indexTmpl, baseTmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// render the templates
	if err := tmpl.ExecuteTemplate(w, "base", credList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
