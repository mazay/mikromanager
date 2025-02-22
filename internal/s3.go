package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"path"
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

// s3BasePath returns the base path for a device's exports in the S3 bucket.
// The returned path is <bucketPath>/exports/<deviceId>.
func s3BasePath(bucketPath string, deviceId string) string {
	return path.Join(bucketPath, "exports", deviceId)
}

// GetS3Session configures the AWS S3 client and returns an error if something fails.
// If the region is not set, it will be determined from the environment. If the
// endpoint is not set, AWS will use the default endpoint for the region. If the
// access key and secret access key are not set, temporary credentials will be
// used.
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

// GetObjects retrieves a list of objects from the S3 bucket with the specified
// prefix. It paginates through the results and returns a slice of S3 objects
// found under the given prefix. If the operation fails, it returns an error.
func (b *S3) GetObjects(prefix string) ([]types.Object, error) {
	var (
		err   error
		items = []types.Object{}
	)

	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(b.Bucket),
		Prefix: aws.String(prefix),
	}

	p := s3.NewListObjectsV2Paginator(b.client, params)

	for p.HasMorePages() {
		page, err := p.NextPage(context.TODO())
		if err != nil {
			return items, err
		}

		// add page contents to the list
		items = append(items, page.Contents...)
	}

	return items, err
}

// GetExports returns a list of all exports for a given device ID in the S3
// bucket. It returns an error if the listing fails.
func (b *S3) GetExports(deviceId string) ([]*Export, error) {
	var (
		err   error
		items = []*Export{}
	)

	prefix := s3BasePath(b.BucketPath, deviceId)
	objects, err := b.GetObjects(prefix)
	if err != nil {
		return items, err
	}

	// convert objects to Export
	for _, s3obj := range objects {
		items = append(items, &Export{
			Key:          *s3obj.Key,
			DeviceId:     deviceId,
			LastModified: s3obj.LastModified,
			ETag:         *s3obj.ETag,
			Size:         s3obj.Size,
		})
	}

	return items, err
}

// GetObjectAttributes returns the attributes of the specified S3 object.
// The object is identified by its key. The returned object contains the ETag,
// Checksum, ObjectParts, StorageClass, and ObjectSize of the object.
// It returns an error if the retrieval fails.
func (b *S3) GetObjectAttributes(key string) (*s3.GetObjectAttributesOutput, error) {
	params := &s3.GetObjectAttributesInput{
		Bucket: aws.String(b.Bucket),
		Key:    aws.String(key),
		ObjectAttributes: []types.ObjectAttributes{
			types.ObjectAttributesEtag,
			types.ObjectAttributesChecksum,
			types.ObjectAttributesObjectParts,
			types.ObjectAttributesStorageClass,
			types.ObjectAttributesObjectSize,
		},
	}

	return b.client.GetObjectAttributes(context.TODO(), params)
}

// GetExportAttributes retrieves the attributes of an export object from the S3
// bucket and returns them as an Export object. The object is identified by its
// key. The returned object contains the ETag, LastModified, and Size of the
// object. It returns an error if the retrieval fails.
func (b *S3) GetExportAttributes(key string) (*Export, error) {
	var (
		err    error
		export = &Export{}
	)

	resp, err := b.GetObjectAttributes(key)

	export.Key = key
	export.LastModified = resp.LastModified
	export.ETag = *resp.ETag
	export.Size = resp.ObjectSize

	return export, err
}

// UploadFile uploads a file to the S3 bucket using the provided S3 key.
// The file data is split into parts of 5 MB each for upload, and the
// upload is performed concurrently with a concurrency level of 5.
// The StorageClass specified in the S3 struct is used for the upload.
// It returns an error if the upload fails.
func (b *S3) UploadFile(s3Key string, data []byte) error {
	uploader := manager.NewUploader(b.client, func(u *manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024 // 5 MB per part
		u.Concurrency = 5
	})

	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:       aws.String(b.Bucket),
		Key:          aws.String(s3Key),
		Body:         bytes.NewReader(data),
		StorageClass: types.StorageClass(b.StorageClass),
	})

	return err
}

// UploadExport uploads an export for a given device ID to the S3 bucket.
// The export data is stored in a file with a name like <unix-timestamp>.rsc in
// a directory like <bucketPath>/exports/<deviceId>. The StorageClass specified
// in the S3 struct is used for the upload. It returns an error if the upload
// fails.
func (b *S3) UploadExport(deviceId string, data []byte) error {
	s3Key := path.Join(s3BasePath(b.BucketPath, deviceId), fmt.Sprintf("%d.rsc", time.Now().Unix()))
	return b.UploadFile(s3Key, data)
}

// GetFile downloads a file from the S3 bucket using the provided S3 key and size.
// It splits the download into parts of 5 MB each and performs the download concurrently
// with a concurrency level of 5. The function returns the downloaded file contents as a
// byte slice and an error if the download fails.
func (b *S3) GetFile(s3Key string, size int64) ([]byte, error) {
	downloader := manager.NewDownloader(b.client, func(d *manager.Downloader) {
		d.PartSize = 5 * 1024 * 1024 // 5 MB per part
		d.Concurrency = 5
	})

	buf := make([]byte, int(size))
	w := manager.NewWriteAtBuffer(buf)
	_, err := downloader.Download(context.TODO(), w, &s3.GetObjectInput{
		Bucket: aws.String(b.Bucket),
		Key:    aws.String(s3Key),
	})

	return buf, err
}

// DeleteFile removes an object from the S3 bucket specified by the S3 key.
// It returns an error if the deletion fails. The S3 key should be a valid
// path to the object within the bucket.
func (b *S3) DeleteFile(s3Key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(b.Bucket),
		Key:    aws.String(s3Key),
	}

	// Delete the object
	_, err := b.client.DeleteObject(context.TODO(), input)

	return err
}

// DeleteExports removes a list of exports from the S3 bucket specified by the S3 keys
// in the list of Export objects. If the list is empty, the function simply returns nil.
// It returns an error if the deletion fails for any of the objects.
func (b *S3) DeleteExports(exports []*Export) error {
	if len(exports) == 0 {
		return nil
	}

	// gererate a set of ObjectIdentifiers
	objects := make([]types.ObjectIdentifier, len(exports))
	for i, e := range exports {
		objects[i] = types.ObjectIdentifier{
			Key: aws.String(e.Key),
		}
	}

	input := s3.DeleteObjectsInput{
		Bucket: aws.String(b.Bucket),
		Delete: &types.Delete{
			Objects: objects,
			Quiet:   aws.Bool(true),
		},
	}
	delOut, err := b.client.DeleteObjects(context.TODO(), &input)
	if err != nil || len(delOut.Errors) > 0 {
		if err != nil {
			var noBucket *types.NoSuchBucket
			if errors.As(err, &noBucket) {
				err = noBucket
			}
		} else if len(delOut.Errors) > 0 {
			for _, outErr := range delOut.Errors {
				log.Printf("%s: %s\n", *outErr.Key, *outErr.Message)
			}
			err = fmt.Errorf("%s", *delOut.Errors[0].Message)
		}
	} else {
		for _, delObjs := range delOut.Deleted {
			err = s3.NewObjectNotExistsWaiter(b.client).Wait(
				context.TODO(), &s3.HeadObjectInput{Bucket: aws.String(b.Bucket), Key: delObjs.Key}, time.Minute)
		}
	}

	return err
}
