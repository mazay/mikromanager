package internal

import (
	"path/filepath"
	"time"
)

type Export struct {
	Key          string
	DeviceId     string
	LastModified *time.Time
	ETag         string
	Size         *int64
}

// GetDeviceId returns the device ID based on the export's key.
// The device ID is inferred from the directory name of the export's key.
func (e *Export) GetDeviceId() string {
	dir := filepath.Dir(e.Key)
	return filepath.Base(dir)
}

// GetBody downloads the export from the S3 bucket and returns its contents as a byte slice.
// It returns an error if the download fails.
func (e *Export) GetBody(s3 *S3) ([]byte, error) {
	return s3.GetFile(e.Key, *e.Size)
}
