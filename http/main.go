package http

import (
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/mazay/mikromanager/db"
	"github.com/sirupsen/logrus"
)

var (
	baseTmpl = path.Join("templates", "base.html")
)

type dynamicHandler struct {
	db            *db.DB
	encryptionKey string
	logger        *logrus.Entry
	backupPath    string
}

func handlerWrapper(fn http.HandlerFunc, logger *logrus.Entry) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Infof("%s - \"%s %s %s\" %s", req.RemoteAddr, req.Method,
			req.URL.Path, req.Proto, strconv.FormatInt(req.ContentLength, 10))

		res.Header().Set("Server", "mikromanager")

		fn(res, req)
	}
}

func HttpServer(httpPort string, db *db.DB, encryptionKey string, backupPath string, logger *logrus.Entry) {
	logger.Infof("starting http server on port %s", httpPort)
	dh := dynamicHandler{db: db, encryptionKey: encryptionKey, logger: logger, backupPath: backupPath}
	static := http.FileServer(http.Dir("./static"))
	backups := http.FileServer(http.Dir(backupPath))
	http.HandleFunc("/", handlerWrapper(dh.getDevices, logger))
	http.HandleFunc("/details", handlerWrapper(dh.getDevice, logger))
	http.HandleFunc("/edit", handlerWrapper(dh.editDevice, logger))
	http.HandleFunc("/delete", handlerWrapper(dh.deleteDevice, logger))
	http.HandleFunc("/credentials", handlerWrapper(dh.getCredentials, logger))
	http.HandleFunc("/credentials/edit", handlerWrapper(dh.editCredentials, logger))
	http.HandleFunc("/credentials/delete", handlerWrapper(dh.deleteCredentials, logger))
	http.Handle("/static/", http.StripPrefix("/static/", static))
	http.Handle("/backups/", http.StripPrefix("/backups/", backups))
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}
