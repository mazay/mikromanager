package http

import (
	"net/http"

	"github.com/mazay/mikromanager/utils"
)

type userForm struct {
	Id                string
	Username          string
	EncryptedPassword string
	Msg               string
}

type usersData struct {
	Count       int
	Users       []*utils.User
	Pagination  *Pagination
	CurrentPage int
}

func (uf *userForm) formFillIn(user *utils.User) {
	uf.Id = user.Id
	uf.Username = user.Username
	uf.EncryptedPassword = user.EncryptedPassword
}

func (dh *dynamicHandler) editUser(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		formErr   error
		data      = &userForm{}
		user      = &utils.User{}
		templates = []string{userFormTmpl, baseTmpl}
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
		username := r.PostForm.Get("username")
		encryptedPw, err := utils.EncryptString(r.PostForm.Get("password"), dh.encryptionKey)
		if err != nil {
			dh.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user.Id = id
		user.Username = username
		user.EncryptedPassword = encryptedPw

		if id == "" {
			// "id" is unset - create new user
			formErr = user.Create(dh.db)
			if formErr != nil {
				data.Msg = formErr.Error()
			}
		} else {
			// "id" is set - update existing user
			formErr = user.GetById(dh.db)
			if formErr != nil {
				data.Msg = err.Error()
			}
			user.Username = username
			user.EncryptedPassword = encryptedPw
			formErr = user.Update(dh.db)
			if formErr != nil {
				data.Msg = err.Error()
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
			u := &utils.User{Id: id}
			err = u.GetById(dh.db)
			if err != nil {
				data.Msg = err.Error()
			} else {
				data.formFillIn(u)
			}
		}
	}

	dh.renderTemplate(w, templates, data)
}

func (dh *dynamicHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		u          = &utils.User{}
		data       = &usersData{}
		pagination = &Pagination{}
		templates  = []string{usersTmpl, paginationTmpl, baseTmpl}
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

	userList, err := u.GetAll(dh.db)
	if err != nil {
		dh.logger.Error(err.Error())
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

	dh.renderTemplate(w, templates, data)
}

func (dh *dynamicHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		u   = &utils.User{}
		id  = r.URL.Query().Get("id")
	)

	_, err = dh.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if id == "" {
		http.Error(w, "Something went wrong, no user ID provided", http.StatusInternalServerError)
		return
	}

	u.Id = id

	err = u.Delete(dh.db)
	if err != nil {
		dh.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusFound)
}
