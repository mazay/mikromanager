package http

import (
	"net/http"
	"strconv"

	"github.com/mazay/mikromanager/db"
)

type exportRetentionPolicyForm struct {
	Id     string
	Name   string
	Hourly int64
	Daily  int64
	Weekly int64
	Msg    string
}

func (erp *exportRetentionPolicyForm) formFillIn(policy *db.ExportsRetentionPolicy) {
	erp.Id = policy.Id
	erp.Name = policy.Name
	erp.Hourly = policy.Hourly
	erp.Daily = policy.Daily
	erp.Weekly = policy.Weekly
}

func (c *HttpConfig) editExportRetentionPolicy(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		data      = &exportRetentionPolicyForm{}
		erp       = &db.ExportsRetentionPolicy{Name: "Default"}
		templates = []string{erpTmpl, baseTmpl}
	)

	_, err = c.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	err = erp.GetDefault(c.Db)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

		hourly, err := strconv.ParseInt(r.PostForm.Get("hourly"), 10, 64)
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		erp.Hourly = hourly

		daily, err := strconv.ParseInt(r.PostForm.Get("daily"), 10, 64)
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		erp.Daily = daily

		weekly, err := strconv.ParseInt(r.PostForm.Get("weekly"), 10, 64)
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		erp.Weekly = weekly

		err = erp.Update(c.Db)
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	data.formFillIn(erp)
	c.renderTemplate(w, templates, data)
}
