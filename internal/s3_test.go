package internal

import (
	"reflect"
	"testing"
)

func TestS3BasePath(t *testing.T) {
	tests := []struct {
		name       string
		bucketPath string
		deviceId   string
		want       string
	}{
		{"Empty bucketPath and deviceId", "", "", "exports"},
		{"Empty bucketPath", "", "device-id", "exports/device-id"},
		{"Empty deviceId", "bucket", "", "bucket/exports"},
		{"Both bucketPath and deviceId", "bucket", "device-id", "bucket/exports/device-id"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s3BasePath(tt.bucketPath, tt.deviceId); got != tt.want {
				t.Errorf("s3BasePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetS3Session(t *testing.T) {
	tests := []struct {
		name            string
		region          string
		endpoint        string
		accessKey       string
		secretAccessKey string
		wantErr         bool
	}{
		{"Empty region and endpoint", "", "", "", "", false},
		{"Region but no endpoint", "us-east-1", "", "", "", false},
		{"Endpoint but no region", "", "https://s3.amazonaws.com", "", "", false},
		{"Both region and endpoint", "us-east-1", "https://s3.amazonaws.com", "", "", false},
		{"URL without scheme", "", " invalid-url", "", "", false},
		{"Access key and secret access key", "us-east-1", "", "access-key", "secret-access-key", false},
		{"No access key and secret access key", "us-east-1", "", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &S3{
				Region:          tt.region,
				Endpoint:        tt.endpoint,
				AccessKey:       tt.accessKey,
				SecretAccessKey: tt.secretAccessKey,
			}

			err := b.GetS3Session()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetS3Session() error = %v, wantErr %v", err, tt.wantErr)
			}
			responseType := reflect.TypeOf(b.client).String()
			if responseType != "*s3.Client" {
				t.Error("Expected type *s3.Client, got", responseType)
			}
		})
	}
}
