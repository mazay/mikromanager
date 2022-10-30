package http

import (
	"net/http"

	"github.com/mazay/mikromanager/db"
	"github.com/sirupsen/logrus"
)

func HttpServer(httpPort string, db *db.DB, encryptionKey string, backupPath string, logger *logrus.Entry) {
	logger.Infof("starting http server on port %s", httpPort)
	dh := &dynamicHandler{db: db, encryptionKey: encryptionKey, logger: logger, backupPath: backupPath}
	static := http.FileServer(http.Dir("./static"))
	backups := http.FileServer(http.Dir(backupPath))
	http.HandleFunc("/healthz", handlerWrapper(dh.healthz, logger))
	http.HandleFunc("/", handlerWrapper(dh.getDevices, logger))
	http.HandleFunc("/details", handlerWrapper(dh.getDevice, logger))
	http.HandleFunc("/edit", handlerWrapper(dh.editDevice, logger))
	http.HandleFunc("/delete", handlerWrapper(dh.deleteDevice, logger))
	http.HandleFunc("/credentials", handlerWrapper(dh.getCredentials, logger))
	http.HandleFunc("/credentials/edit", handlerWrapper(dh.editCredentials, logger))
	http.HandleFunc("/credentials/delete", handlerWrapper(dh.deleteCredentials, logger))
	http.HandleFunc("/erp", handlerWrapper(dh.editExportRetentionPolicy, logger))
	http.HandleFunc("/exports", handlerWrapper(dh.getExports, logger))
	http.Handle("/static/", http.StripPrefix("/static/", static))
	http.Handle("/backups/", http.StripPrefix("/backups/", backups))
	logger.Fatal(http.ListenAndServe(":"+httpPort, nil))
}
