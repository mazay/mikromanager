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
	db            *db.DB
	encryptionKey string
}

func handlerWrapper(fn http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("%s - \"%s %s %s\" %s", req.RemoteAddr, req.Method,
			req.URL.Path, req.Proto, strconv.FormatInt(req.ContentLength, 10))

		res.Header().Set("Server", "mikromanager")

		fn(res, req)
	}
}

func HttpServer(httpPort string, db *db.DB, encryptionKey string) {
	log.Printf("starting http server on port %s", httpPort)
	dh := dynamicHandler{db: db, encryptionKey: encryptionKey}
	fs := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/", handlerWrapper(dh.getDevices))
	http.HandleFunc("/edit", handlerWrapper(dh.editDevice))
	http.HandleFunc("/delete", handlerWrapper(dh.deleteDevice))
	http.HandleFunc("/credentials", handlerWrapper(dh.getCredentials))
	http.HandleFunc("/credentials/edit", handlerWrapper(dh.editCredentials))
	http.HandleFunc("/credentials/delete", handlerWrapper(dh.deleteCredentials))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}
