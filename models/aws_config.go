package models

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	zone              = "zone"
	bucket            = "bucket"
	path              = "path"
	contentType       = "contentType"
	distributionID    = "distributionID"
	invalidationBatch = "invalidationBatch"
)

//AwsStorageConfig defined the storage type we will use in AWS and the location of the bucket
type AwsStorageConfig interface {
	Zone() string
	Bucket() string
	ContentType() string
	Path() string
}

//S3Config contains the host and port for the dao
type S3Config struct {
	zone        string
	bucket      string
	contentType string
	path        string
}

//Zone retrieve zone
func (c *S3Config) Zone() string {
	return c.zone
}

//Bucket retrieve bucket
func (c *S3Config) Bucket() string {
	return c.bucket
}

//ContentType retrieve contentType
func (c *S3Config) ContentType() string {
	return c.contentType
}

//Path retrieve path
func (c *S3Config) Path() string {
	return c.path
}

//S3Client contains an Aws Session and a  that is used to interact with the db
type S3Client struct {
	Session *session.Session
	Config  AwsStorageConfig
}

//NewS3Config creates an S3 config based on environment vars passed to the Lambda
func NewS3Config() *S3Config {
	zone := os.Getenv(zone)
	bucket := os.Getenv(bucket)
	contentType := os.Getenv(contentType)
	path := resolvePath(os.Getenv(path))

	if zone == "" || bucket == "" || contentType == "" {
		return nil
	}

	return &S3Config{
		zone,
		bucket,
		contentType,
		path,
	}
}

//resolvePath validate and resolve path location
func resolvePath(path string) string {
	if len(strings.TrimSpace(path)) == 0 {
		return path
	}
	cleanedPath := strings.Replace(filepath.Clean(path), "\\", "/", -1)
	if cleanedPath[:1] == "/" {
		cleanedPath = cleanedPath[1:]
	}
	if cleanedPath[len(cleanedPath)-1:] != "/" {
		cleanedPath = cleanedPath + "/"
	}
	return cleanedPath
}

//NewS3Client accepts configuration type that can setup and return a S3 session to the zone of interest
func NewS3Client(configuration AwsStorageConfig) (*S3Client, error) {
	if configuration == nil {
		return nil, errors.New("nil config passed, failed to initialize config for the AWS connections")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(configuration.Zone())},
	)
	if err != nil {
		return nil, err
	}

	return &S3Client{
		Session: sess,
		Config:  configuration,
	}, nil
}
