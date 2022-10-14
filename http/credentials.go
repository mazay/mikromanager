package http

import (
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/mazay/mikromanager/utils"
)

type credentialsForm struct {
	Alias    string
	Username string
	Msg      string
}

func (dh *dynamicHandler) getCredentials(w http.ResponseWriter, r *http.Request) {
	var credsTmpl = path.Join("templates", "credentials.html")
	var c = &utils.Credentials{}

	// fetch devices
	credList, err := c.GetAll(dh.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// load templates
	tmpl, err := template.New("").ParseFiles(credsTmpl, baseTmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// render the templates
	if err := tmpl.ExecuteTemplate(w, "base", credList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (dh *dynamicHandler) addCredentials(w http.ResponseWriter, r *http.Request) {
	var (
		credsTmpl = path.Join("templates", "credentials-form.html")
		data      = &credentialsForm{}
	)

	if r.Method == "POST" {
		r.ParseForm()
		encryptedPw, err := utils.EncryptString(r.PostForm.Get("password"), dh.encryptionKey)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		creds := &utils.Credentials{
			Alias:             r.PostForm.Get("alias"),
			Username:          r.PostForm.Get("username"),
			EncryptedPassword: encryptedPw,
		}
		credsErr := creds.Create(dh.db)
		if credsErr != nil {
			data.Alias = r.PostForm.Get("alias")
			data.Username = r.PostForm.Get("username")
			data.Msg = credsErr.Error()
		} else {
			http.Redirect(w, r, "/credentials", 302)
			return
		}
	}

	// load templates
	tmpl, err := template.New("").ParseFiles(credsTmpl, baseTmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// render the templates
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
