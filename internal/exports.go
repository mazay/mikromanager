package internal

import (
	"time"
)

type Export struct {
	Key          string
	DeviceId     string
	LastModified *time.Time
	ETag         string
	Size         *int64
}
