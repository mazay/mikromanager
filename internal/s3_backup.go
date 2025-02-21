package internal

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3 struct {
	Bucket          string
	BucketPath      string
	Endpoint        string
	Region          string
	StorageClass    string
	AccessKey       string
	SecretAccessKey string
	OpsRetries      int
	client          *s3.Client
}

func s3BasePath(bucketPath string, deviceId string) string {
	return path.Join(bucketPath, "exports", deviceId)
}

func (b *S3) GetS3Session() error {
	var (
		err      error
		endpoint *url.URL
	)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(b.Region),
	)
	if err != nil {
		return err
	}

	if b.Endpoint != "" {
		endpoint, err = url.Parse(b.Endpoint)
		if err != nil {
			return err
		}
		if endpoint.Scheme == "" {
			endpoint.Scheme = "https"
		}
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		// override endpoint if needed
		if b.Endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint.String())
		}
		// set static credentials if needed
		if b.AccessKey != "" && b.SecretAccessKey != "" {
			o.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(b.AccessKey, b.SecretAccessKey, ""))
		}
		o.RetryMaxAttempts = b.OpsRetries
	})

	b.client = s3Client

	return err
}

func (b *S3) GetAwsS3ItemMap(deviceId string) (map[string]string, error) {
	var (
		err   error
		iter  int
		items = make(map[string]string)
	)

	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(b.Bucket),
		Prefix: aws.String(s3BasePath(b.BucketPath, deviceId)),
	}

	p := s3.NewListObjectsV2Paginator(b.client, params)

	for p.HasMorePages() {
		iter++

		page, err := p.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		// Log the objects found
		for _, s3obj := range page.Contents {
			if string(s3obj.StorageClass) != b.StorageClass {
				items[*aws.String(*s3obj.Key)] = "none"
			} else {
				// Update metrics
				items[*aws.String(*s3obj.Key)] = strings.Trim(*(s3obj.ETag), "\"")
			}
		}
	}

	return items, err
}

func (b *S3) UploadFile(deviceId string, data []byte) error {
	var err error

	s3Key := path.Join(s3BasePath(b.BucketPath, deviceId), fmt.Sprintf("%d.rsc", time.Now().Unix()))
	uploader := manager.NewUploader(b.client, func(u *manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		u.Concurrency = 5
	})

	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:       aws.String(b.Bucket),
		Key:          aws.String(s3Key),
		Body:         bytes.NewReader(data),
		StorageClass: types.StorageClass(b.StorageClass),
	})

	return err
}

func (b *S3) DeleteFile(s3Key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(b.Bucket),
		Key:    aws.String(s3Key),
	}

	// Delete the object
	_, err := b.client.DeleteObject(context.TODO(), input)

	return err
}
