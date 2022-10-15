package main

import (
	"errors"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	ApiPollers    int    `yaml:"apiPollers"`
	DbPath        string `yaml:"dbPath"`
	EncryptionKey string `yaml:"encryptionKey"`
	LogLevel      string `yaml:"logLevel"`
}

func configProcessError(err error) {
	logger.Error(err)
	osExit(2)
}

func (cfg *Config) setDefaults() {
	if cfg.ApiPollers == 0 {
		// there should be at least 1 poller
		cfg.ApiPollers = 1
	}
	if cfg.DbPath == "" {
		cfg.DbPath = "database.clover"
	}
	if cfg.EncryptionKey == "" {
		configProcessError(errors.New("the encryptionKey should be set"))
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
}

func readConfigFile(path string) *Config {
	var inInterface Config
	f, err := os.ReadFile(path)
	if err != nil {
		configProcessError(err)
	}

	err = yaml.Unmarshal(f, &inInterface)
	if err != nil {
		configProcessError(err)
	} else {
		inInterface.setDefaults()
	}
	return &inInterface
}