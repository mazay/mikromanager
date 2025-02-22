package http

import (
	"net/http"

	"github.com/mazay/mikromanager/utils"
)

type credentialsForm struct {
	Id       string
	Alias    string
	Username string
	Msg      string
}

type credentialsData struct {
	Count       int
	Credentials []*utils.Credentials
	Pagination  *Pagination
	CurrentPage int
}

func (cf *credentialsForm) formFillIn(creds *utils.Credentials) {
	cf.Id = creds.Id
	cf.Alias = creds.Alias
	cf.Username = creds.Username
}

func (dh *dynamicHandler) getCredentials(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		c          = &utils.Credentials{}
		data       = &credentialsData{}
		pagination = &Pagination{}
		templates  = []string{credsTmpl, paginationTmpl, baseTmpl}
	)

	_, err = dh.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	pageId, perPage, err := getPagionationParams(r.URL)
	if err != nil {
		dh.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fetch devices
	credList, err := c.GetAll(dh.db)
	if err != nil {
		dh.logger.Error(err.Error())
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

	dh.renderTemplate(w, templates, data)
}

func (dh *dynamicHandler) editCredentials(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		credsErr  error
		data      = &credentialsForm{}
		templates = []string{credsFormTmpl, baseTmpl}
	)

	_, err = dh.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method == "POST" {
		// parse the form
		err = r.ParseForm()
		if err != nil {
			dh.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := r.PostForm.Get("idInput")
		alias := r.PostForm.Get("alias")
		username := r.PostForm.Get("username")
		encryptedPw, err := utils.EncryptString(r.PostForm.Get("password"), dh.encryptionKey)

		if err != nil {
			dh.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		creds := &utils.Credentials{
			Id:                id,
			Alias:             alias,
			Username:          username,
			EncryptedPassword: encryptedPw,
		}

		if id == "" {
			// "id" is unset - create new credentials
			credsErr = creds.Create(dh.db)
		} else {
			// "id" is set - update existing credentials
			err = creds.GetById(dh.db)
			if err != nil {
				dh.logger.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			creds.Alias = alias
			creds.Username = username
			creds.EncryptedPassword = encryptedPw
			credsErr = creds.Update(dh.db)
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
			c := &utils.Credentials{}
			c.Id = id
			err := c.GetById(dh.db)
			if err != nil {
				data.Msg = err.Error()
			} else {
				data.formFillIn(c)
			}
		}
	}

	dh.renderTemplate(w, templates, data)
}

func (dh *dynamicHandler) deleteCredentials(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		c   = &utils.Credentials{}
	)

	_, err = dh.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	c.Id = r.URL.Query().Get("id")

	err = c.Delete(dh.db)
	if err != nil {
		dh.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/credentials", http.StatusFound)
}
