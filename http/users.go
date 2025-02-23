package http

import (
	"net/http"

	"github.com/mazay/mikromanager/db"
)

type userForm struct {
	Id                string
	Username          string
	EncryptedPassword string
	Msg               string
}

type usersData struct {
	Count       int
	Users       []*db.User
	Pagination  *Pagination
	CurrentPage int
}

func (uf *userForm) formFillIn(user *db.User) {
	uf.Id = user.ID
	uf.Username = user.Username
	uf.EncryptedPassword = user.EncryptedPassword
}

func (c *HttpConfig) editUser(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		formErr   error
		data      = &userForm{}
		user      = &db.User{}
		templates = []string{userFormTmpl, baseTmpl}
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
		username := r.PostForm.Get("username")
		encryptedPw, err := db.EncryptString(r.PostForm.Get("password"), c.EncryptionKey)
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user.ID = id
		user.Username = username
		user.EncryptedPassword = encryptedPw

		if id == "" {
			// "id" is unset - create new user
			formErr = user.Create(c.Db)
			if formErr != nil {
				data.Msg = formErr.Error()
			}
		} else {
			// "id" is set - update existing user
			formErr = user.GetById(c.Db)
			if formErr != nil {
				data.Msg = formErr.Error()
			}
			user.Username = username
			user.EncryptedPassword = encryptedPw
			formErr = user.Update(c.Db)
			if formErr != nil {
				data.Msg = formErr.Error()
			}
		}

		if formErr != nil {
			// return data with errors if validation failed
			data.Id = id
			data.formFillIn(user)
		} else {
			http.Redirect(w, r, "/users", http.StatusFound)
			return
		}
	} else {
		// fill in the form if "id" GET parameter set
		id := r.URL.Query().Get("id")
		if id != "" {
			u := &db.User{}
			u.ID = id
			err = u.GetById(c.Db)
			if err != nil {
				data.Msg = err.Error()
			} else {
				data.formFillIn(u)
			}
		}
	}

	c.renderTemplate(w, templates, data)
}

func (c *HttpConfig) getUsers(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		u          = &db.User{}
		data       = &usersData{}
		pagination = &Pagination{}
		templates  = []string{usersTmpl, paginationTmpl, baseTmpl}
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

	userList, err := u.GetAll(c.Db)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data.Count = len(userList)
	if data.Count > 0 {
		chunkedUsers := chunkSliceOfObjects(userList, perPage)
		pagination.paginate(*r.URL, pageId, len(chunkedUsers))

		if pageId-1 >= len(chunkedUsers) {
			pageId = len(chunkedUsers)
		}
		data.Pagination = pagination
		data.CurrentPage = pageId
		data.Users = chunkedUsers[pageId-1]
	}

	c.renderTemplate(w, templates, data)
}

func (c *HttpConfig) deleteUser(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		u   = &db.User{}
		id  = r.URL.Query().Get("id")
	)

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if id == "" {
		http.Error(w, "Something went wrong, no user ID provided", http.StatusInternalServerError)
		return
	}

	u.ID = id

	err = u.Delete(c.Db)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusFound)
}
