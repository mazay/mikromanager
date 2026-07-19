package http

import (
	"net/http"

	"github.com/mazay/mikromanager/db"
)

type snippetForm struct {
	Id      string
	Name    string
	Content string
	Msg     string
}

func (f *snippetForm) formFillIn(cs *db.ConfigurationSnippet) {
	f.Id = cs.Id
	f.Name = cs.Name
	f.Content = cs.Content
}

func (c *HttpConfig) editSnippet(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		csErr     error
		data      = &snippetForm{}
		templates = []string{csFormTmpl, baseTmpl}
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
		name := r.PostForm.Get("name")
		content := r.PostForm.Get("content")

		cs := &db.ConfigurationSnippet{
			Name:    name,
			Content: content,
		}

		if id == "" {
			// "id" is unset - create new snippet
			csErr = cs.Create(c.Db)
		} else {
			// "id" is set - update existing snippet
			cs.Id = id
			err = cs.GetById(c.Db)
			if err != nil {
				c.Logger.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			cs.Name = name
			cs.Content = content
			csErr = cs.Update(c.Db)
		}

		if csErr != nil {
			// return data with errors if validation failed
			data.formFillIn(cs)
			data.Msg = csErr.Error()
		} else {
			http.Redirect(w, r, "/config-snippets", http.StatusFound)
			return
		}
	} else {
		// fill in the form if "id" GET parameter set
		id := r.URL.Query().Get("id")
		if id != "" {
			cs := &db.ConfigurationSnippet{}
			cs.Id = id
			err = cs.GetById(c.Db)
			if err != nil {
				data.Msg = err.Error()
			} else {
				data.formFillIn(cs)
			}
		}
	}

	c.renderTemplate(w, templates, data)
}
