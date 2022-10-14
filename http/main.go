package http

import (
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/mazay/mikromanager/db"
)

var (
	baseTmpl = path.Join("templates", "base.html")
)

type dynamicHandler struct {
	db *db.DB
}

func handlerWrapper(fn http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("%s - \"%s %s %s\" %s", req.RemoteAddr, req.Method,
			req.URL.Path, req.Proto, strconv.FormatInt(req.ContentLength, 10))

		res.Header().Set("Server", "mikromanager")

		fn(res, req)
	}
}

func HttpServer(httpPort string, db *db.DB) {
	log.Printf("starting http server on port %s", httpPort)
	dh := dynamicHandler{db: db}
	fs := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/", handlerWrapper(dh.getDevices))
	http.HandleFunc("/credentials", handlerWrapper(dh.getCredentials))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}
