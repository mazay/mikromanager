package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func LoggerInitError(err error) {
	logger.Error(err)
	osExit(3)
}

func setLogLevel(level string) {
	logLevels := map[string]logrus.Level{
		"trace": logrus.TraceLevel,
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
		"fatal": logrus.FatalLevel,
		"panic": logrus.PanicLevel,
	}

	if level == "" {
		level = "info"
	}

	if LogLevel, ok := logLevels[level]; ok {
		log.SetLevel(LogLevel)

		if level == "trace" || level == "debug" {
			log.SetReportCaller(true)
		}
	} else {
		LoggerInitError(fmt.Errorf("Log level definition not found for '%s'", level))
	}
}

func initLogger(config *Config) {
	log.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	})

	log.SetOutput(os.Stdout)

	setLogLevel(config.LogLevel)
}
