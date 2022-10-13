package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	database "github.com/mazay/mikromanager/db"
)

type dynamicHandler struct {
	db *database.DB
}

func handlerWrapper(fn http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("%s - \"%s %s %s\" %s", req.RemoteAddr, req.Method,
			req.URL.Path, req.Proto, strconv.FormatInt(req.ContentLength, 10))

		res.Header().Set("Server", "mikromanager")

		fn(res, req)
	}
}

func (dh *dynamicHandler) getRoot(w http.ResponseWriter, r *http.Request) {
	devices, err := dh.db.FindAll("devices")
	if err != nil {
		io.WriteString(w, err.Error())
	}
	data, err := json.Marshal(devices)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(w, string(data))
}

func HttpServer(httpPort string, db *database.DB) {
	log.Printf("starting http server on port %s", httpPort)
	dh := dynamicHandler{db: db}
	fs := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/", handlerWrapper(dh.getRoot))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}
