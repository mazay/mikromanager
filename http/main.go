package http

import (
	"net/http"

	"github.com/mazay/mikromanager/db"
	"github.com/mazay/mikromanager/internal"
	"go.uber.org/zap"
)

type HttpConfig struct {
	Port          string
	Db            *db.DB
	EncryptionKey string
	Logger        *zap.Logger
	BackupPath    string
	S3            *internal.S3
}

func (c *HttpConfig) HttpServer() {
	c.Logger.Info("starting http server", zap.String("port", c.Port))
	static := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/healthz", handlerWrapper(c.healthz, c.Logger))
	http.HandleFunc("/login", handlerWrapper(c.login, c.Logger))
	http.HandleFunc("/logout", handlerWrapper(c.logout, c.Logger))
	http.HandleFunc("/users", handlerWrapper(c.getUsers, c.Logger))
	http.HandleFunc("/user/edit", handlerWrapper(c.editUser, c.Logger))
	http.HandleFunc("/user/delete", handlerWrapper(c.deleteUser, c.Logger))
	http.HandleFunc("/", handlerWrapper(c.getDevices, c.Logger))
	http.HandleFunc("/details", handlerWrapper(c.getDevice, c.Logger))
	http.HandleFunc("/edit", handlerWrapper(c.editDevice, c.Logger))
	http.HandleFunc("/delete", handlerWrapper(c.deleteDevice, c.Logger))
	http.HandleFunc("/credentials", handlerWrapper(c.getCredentials, c.Logger))
	http.HandleFunc("/credentials/edit", handlerWrapper(c.editCredentials, c.Logger))
	http.HandleFunc("/credentials/delete", handlerWrapper(c.deleteCredentials, c.Logger))
	http.HandleFunc("/erp", handlerWrapper(c.editExportRetentionPolicy, c.Logger))
	http.HandleFunc("/exports", handlerWrapper(c.getExports, c.Logger))
	http.HandleFunc("/export", handlerWrapper(c.getExport, c.Logger))
	http.HandleFunc("/export/download", handlerWrapper(c.downloadExport, c.Logger))
	http.HandleFunc("/device/groups", handlerWrapper(c.getDeviceGroups, c.Logger))
	http.HandleFunc("/device/group/edit", handlerWrapper(c.editDeviceGroup, c.Logger))
	http.HandleFunc("/device/group", handlerWrapper(c.getDeviceGroup, c.Logger))
	http.HandleFunc("/device/group/delete", handlerWrapper(c.deleteDeviceGroup, c.Logger))
	http.Handle("/static/", http.StripPrefix("/static/", static))
	c.Logger.Fatal(http.ListenAndServe(":"+c.Port, nil).Error())
}
