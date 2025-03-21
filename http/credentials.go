package http

import (
	"net/http"

	"github.com/mazay/mikromanager/db"
)

type credentialsForm struct {
	Id       string
	Alias    string
	Username string
	Msg      string
}

type credentialsData struct {
	Count       int
	Credentials []*db.Credentials
	Pagination  *Pagination
	CurrentPage int
}

func (cf *credentialsForm) formFillIn(creds *db.Credentials) {
	cf.Id = creds.Id
	cf.Alias = creds.Alias
	cf.Username = creds.Username
}

func (c *HttpConfig) getCredentials(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		creds      = &db.Credentials{}
		data       = &credentialsData{}
		pagination = &Pagination{}
		templates  = []string{credsTmpl, paginationTmpl, baseTmpl}
	)

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	pageId, perPage, err := getPagionationParams(r.URL)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fetch devices
	credList, err := creds.GetAll(c.Db)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data.Count = len(credList)
	if data.Count > 0 {
		chunkedCreds := chunkSliceOfObjects(credList, perPage)
		pagination.paginate(*r.URL, pageId, len(chunkedCreds))

		if pageId-1 >= len(chunkedCreds) {
			pageId = len(chunkedCreds)
		}

		data.Pagination = pagination
		data.CurrentPage = pageId
		data.Credentials = chunkedCreds[pageId-1]
	}

	c.renderTemplate(w, templates, data)
}

func (c *HttpConfig) editCredentials(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		credsErr  error
		data      = &credentialsForm{}
		templates = []string{credsFormTmpl, baseTmpl}
	)

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method == "POST" {
		// parse the form
		err = r.ParseForm()
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := r.PostForm.Get("idInput")
		alias := r.PostForm.Get("alias")
		username := r.PostForm.Get("username")
		encryptedPw, err := db.EncryptString(r.PostForm.Get("password"), c.EncryptionKey)

		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		creds := &db.Credentials{
			Alias:             alias,
			Username:          username,
			EncryptedPassword: encryptedPw,
		}

		if id == "" {
			// "id" is unset - create new credentials
			credsErr = creds.Create(c.Db)
		} else {
			// "id" is set - update existing credentials
			creds.Id = id
			err = creds.GetById(c.Db)
			if err != nil {
				c.Logger.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			creds.Alias = alias
			creds.Username = username
			creds.EncryptedPassword = encryptedPw
			credsErr = creds.Update(c.Db)
		}

		if credsErr != nil {
			// return data with errors if validation failed
			data.formFillIn(creds)
			data.Msg = credsErr.Error()
		} else {
			http.Redirect(w, r, "/credentials", http.StatusFound)
			return
		}
	} else {
		// fill in the form if "id" GET parameter set
		id := r.URL.Query().Get("id")
		if id != "" {
			creds := &db.Credentials{}
			creds.Id = id
			err = creds.GetById(c.Db)
			if err != nil {
				data.Msg = err.Error()
			} else {
				data.formFillIn(creds)
			}
		}
	}

	c.renderTemplate(w, templates, data)
}

func (c *HttpConfig) deleteCredentials(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		creds = &db.Credentials{}
	)

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	creds.Id = r.URL.Query().Get("id")

	err = creds.Delete(c.Db)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/credentials", http.StatusFound)
}
