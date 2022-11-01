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
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err, pageId, perPage := getPagionationParams(r.URL)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fetch devices
	credList, err := c.GetAll(dh.db)
	if err != nil {
		dh.logger.Error(err)
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
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		// parse the form
		err = r.ParseForm()
		if err != nil {
			dh.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := r.PostForm.Get("idInput")
		alias := r.PostForm.Get("alias")
		username := r.PostForm.Get("username")
		encryptedPw, err := utils.EncryptString(r.PostForm.Get("password"), dh.encryptionKey)

		if err != nil {
			dh.logger.Error(err)
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
				dh.logger.Error(err)
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
			data.Id = id
			data.Alias = alias
			data.Username = username
			data.Msg = credsErr.Error()
		} else {
			http.Redirect(w, r, "/credentials", 302)
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
				data.Id = c.Id
				data.Alias = c.Alias
				data.Username = c.Username
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
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Id = r.URL.Query().Get("id")

	err = c.Delete(dh.db)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/credentials", 302)
}
