package main

import (
	"errors"
	"os"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	ApiPollers           int           `yaml:"apiPollers"`
	BackupPath           string        `yaml:"backupPath"`
	ExportWorkers        int           `yaml:"exportWorkers"`
	DevicePollerInterval time.Duration `yaml:"devicePollerInterval"`
	DeviceExportInterval time.Duration `yaml:"deviceExportInterval"`
	DbPath               string        `yaml:"dbPath"`
	EncryptionKey        string        `yaml:"encryptionKey"`
	LogLevel             string        `yaml:"logLevel"`
	S3Bucket             string        `yaml:"s3Bucket"`
	S3BucketPath         string        `yaml:"s3BucketPath"`
	S3Endpoint           string        `yaml:"s3Endpoint"`
	S3Region             string        `yaml:"s3Region"`
	S3StorageClass       string        `yaml:"s3StorageClass"`
	S3AccessKey          string        `yaml:"s3AccessKey"`
	S3SecretAccessKey    string        `yaml:"s3SecretAccessKey"`
	S3OpsRetries         int           `yaml:"s3OpsRetries"`
}

func configProcessError(err error) {
	logger.Error(err.Error())
	osExit(2)
}

func (cfg *Config) setDefaults() {
	if cfg.ApiPollers == 0 {
		// there should be at least 1 poller
		cfg.ApiPollers = 1
	}
	if cfg.DevicePollerInterval == 0 {
		cfg.DevicePollerInterval = time.Millisecond * 1000 * 300
	}
	if cfg.DeviceExportInterval == 0 {
		cfg.DeviceExportInterval = time.Hour * 1
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
	// S3 sdefaults
	if cfg.S3Region == "" {
		cfg.S3Region = "us-east-1"
	}
	if cfg.S3StorageClass == "" {
		cfg.S3StorageClass = "STANDARD"
	}
	if cfg.S3OpsRetries == 0 {
		cfg.S3OpsRetries = 5
	}
	// check if AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are set
	envAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	if envAccessKey != "" {
		cfg.S3AccessKey = envAccessKey
	}
	envSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if envSecretAccessKey != "" {
		cfg.S3SecretAccessKey = envSecretAccessKey
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
