package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mazay/mikromanager/utils"
)

type loginForm struct {
	Username    string
	Password    string
	UsernameMsg string
	PasswordMsg string
}

func (dh *dynamicHandler) login(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		user    = &utils.User{}
		session = &utils.Session{}
		data    = &loginForm{}
	)

	if r.Method == "POST" {
		// parse the form
		err = r.ParseForm()
		if err != nil {
			dh.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.Username = r.PostForm.Get("username")
		data.Password = r.PostForm.Get("password")

		user.Username = data.Username
		err = user.GetByUsername(dh.db)
		if err != nil {
			dh.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user.Id == "" {
			data.UsernameMsg = fmt.Sprintf("User '%s' not found", data.Username)
			dh.renderTemplate(w, []string{loginTmpl}, data)
			return
		}

		decryptedPw, err := utils.DecryptString(user.EncryptedPassword, dh.encryptionKey)
		if err != nil {
			dh.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if decryptedPw != data.Password {
			data.PasswordMsg = "Incorrect password"
			dh.renderTemplate(w, []string{loginTmpl}, data)
			return
		}

		session.UserId = user.Id
		err = session.Create(dh.db)
		if err != nil {
			dh.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   session.Id,
			Expires: session.ValidThrough,
		})

		http.Redirect(w, r, "/", 302)
		return
	}

	dh.renderTemplate(w, []string{loginTmpl}, data)
}

func (dh *dynamicHandler) logout(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		session = &utils.Session{}
	)

	c, err := r.Cookie("session_token")
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Id = c.Value
	err = session.GetById(dh.db)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = session.Delete(dh.db)
	if err != nil {
		dh.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", 302)
}
