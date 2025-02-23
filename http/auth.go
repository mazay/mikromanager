package http

import (
	"net/http"
	"time"

	"github.com/mazay/mikromanager/db"
)

type loginForm struct {
	Username    string
	Password    string
	UsernameMsg string
	PasswordMsg string
}

func (c *HttpConfig) login(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		user    = &db.User{}
		session = &db.Session{}
		data    = &loginForm{}
	)

	if r.Method == "POST" {
		// parse the form
		err = r.ParseForm()
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.Username = r.PostForm.Get("username")
		data.Password = r.PostForm.Get("password")

		user.Username = data.Username
		err = user.GetByUsername(c.Db)
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		decryptedPw, err := db.DecryptString(user.EncryptedPassword, c.EncryptionKey)
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if decryptedPw != data.Password {
			data.PasswordMsg = "Incorrect password"
			c.renderTemplate(w, []string{loginTmpl}, data)
			return
		}

		session.UserId = user.Id
		err = session.Create(c.Db)
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   session.Id,
			Expires: session.ValidThrough,
		})

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	c.renderTemplate(w, []string{loginTmpl}, data)
}

func (c *HttpConfig) logout(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		session = &db.Session{}
	)

	cookie, err := r.Cookie("session_token")
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Id = cookie.Value
	err = session.GetById(c.Db)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = session.Delete(c.Db)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
