package http

import (
	"net/http"

	"github.com/mazay/mikromanager/db"
	"github.com/mazay/mikromanager/internal"
	"go.uber.org/zap"
)

func HttpServer(httpPort string, db *db.DB, encryptionKey string, backupPath string, logger *zap.Logger, s3 *internal.S3) {
	logger.Info("starting http server", zap.String("port", httpPort))
	dh := &dynamicHandler{db: db, encryptionKey: encryptionKey, logger: logger, backupPath: backupPath, s3: s3}
	static := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/healthz", handlerWrapper(dh.healthz, logger))
	http.HandleFunc("/login", handlerWrapper(dh.login, logger))
	http.HandleFunc("/logout", handlerWrapper(dh.logout, logger))
	http.HandleFunc("/users", handlerWrapper(dh.getUsers, logger))
	http.HandleFunc("/user/edit", handlerWrapper(dh.editUser, logger))
	http.HandleFunc("/user/delete", handlerWrapper(dh.deleteUser, logger))
	http.HandleFunc("/", handlerWrapper(dh.getDevices, logger))
	http.HandleFunc("/details", handlerWrapper(dh.getDevice, logger))
	http.HandleFunc("/edit", handlerWrapper(dh.editDevice, logger))
	http.HandleFunc("/delete", handlerWrapper(dh.deleteDevice, logger))
	http.HandleFunc("/credentials", handlerWrapper(dh.getCredentials, logger))
	http.HandleFunc("/credentials/edit", handlerWrapper(dh.editCredentials, logger))
	http.HandleFunc("/credentials/delete", handlerWrapper(dh.deleteCredentials, logger))
	http.HandleFunc("/erp", handlerWrapper(dh.editExportRetentionPolicy, logger))
	http.HandleFunc("/exports", handlerWrapper(dh.getExports, logger))
	http.HandleFunc("/export", handlerWrapper(dh.getExport, logger))
	http.HandleFunc("/export/download", handlerWrapper(dh.downloadExport, logger))
	http.Handle("/static/", http.StripPrefix("/static/", static))
	logger.Fatal(http.ListenAndServe(":"+httpPort, nil).Error())
}
