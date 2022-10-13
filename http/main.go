package http

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/mazay/mikromanager/db"
	"github.com/mazay/mikromanager/utils"
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

func (dh *dynamicHandler) getRoot(w http.ResponseWriter, r *http.Request) {
	var indexTmpl = path.Join("templates", "index.html")
	var d = &utils.Device{}

	devices, err := dh.db.FindAll("devices")
	if err != nil {
		io.WriteString(w, err.Error())
	}

	deviceList := d.FromListOfMaps(devices)

	tmpl, err := template.New("").ParseFiles(indexTmpl, baseTmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", deviceList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HttpServer(httpPort string, db *db.DB) {
	log.Printf("starting http server on port %s", httpPort)
	dh := dynamicHandler{db: db}
	fs := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/", handlerWrapper(dh.getRoot))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}
