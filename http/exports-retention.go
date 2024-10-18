package http

import (
	"net/http"
	"strconv"

	"github.com/mazay/mikromanager/utils"
)

type exportRetentionPolicyForm struct {
	Id     string
	Name   string
	Hourly int64
	Daily  int64
	Weekly int64
	Msg    string
}

func (erp *exportRetentionPolicyForm) formFillIn(policy *utils.ExportsRetentionPolicy) {
	erp.Id = policy.Id
	erp.Name = policy.Name
	erp.Hourly = policy.Hourly
	erp.Daily = policy.Daily
	erp.Weekly = policy.Weekly
}

func (dh *dynamicHandler) editExportRetentionPolicy(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		data      = &exportRetentionPolicyForm{}
		erp       = &utils.ExportsRetentionPolicy{Name: "Default"}
		templates = []string{erpTmpl, baseTmpl}
	)

	_, err = dh.checkSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	err = erp.GetDefault(dh.db)
	if err != nil {
		dh.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

		hourly, err := strconv.ParseInt(r.PostForm.Get("hourly"), 10, 64)
		if err != nil {
			dh.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		erp.Hourly = hourly

		daily, err := strconv.ParseInt(r.PostForm.Get("daily"), 10, 64)
		if err != nil {
			dh.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		erp.Daily = daily

		weekly, err := strconv.ParseInt(r.PostForm.Get("weekly"), 10, 64)
		if err != nil {
			dh.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		erp.Weekly = weekly

		err = erp.Update(dh.db)
		if err != nil {
			dh.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	data.formFillIn(erp)
	dh.renderTemplate(w, templates, data)
}
