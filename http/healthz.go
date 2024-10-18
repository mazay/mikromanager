package http

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	STATUS string
}

func (dh *dynamicHandler) healthz(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		js       []byte
		response = &healthResponse{}
	)
	w.Header().Set("Content-Type", "application/json")

	response.STATUS = "OK"

	js, err = json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(js)
	if err != nil {
		dh.logger.Error(err.Error())
	}
}
