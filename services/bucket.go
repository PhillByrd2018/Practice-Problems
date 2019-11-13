package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//RunThrough is the main router for AWS services
func RunThrough() {

	fmt.Println("Welcome to AWS run through")
	fmt.Println("Making AWS config")
	svc := CreateAWSConfig()
	fmt.Println("Calling list buckets")
	ListBuckets(svc)
	fmt.Println("Calling create a bucket")
	CreateBucket(svc)
	fmt.Println("Calling insert into a bucket")
	AddItemToBucket(svc, "example-bucket-for-test")
	fmt.Println("Calling list of bucket objects")
	ListobjectsInBucket(svc, "example-bucket-for-test")
	// fmt.Println("Calling list buckets again")
	// ListBuckets(svc)
	GetSingleObject(svc, "example-bucket-for-test", "Solomon")
	bcTest, err := ioutil.ReadFile("services/test.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	sendEmail(string(bcTest))
}

//ListBuckets services used to list buckets in AWS account
func ListBuckets(svc *s3.S3) {

	input := &s3.ListBucketsInput{}

	result, err := svc.ListBuckets(input)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets from AWS: ")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))

		if aws.StringValue(b.Name) == "example-bucket-for-test" {
			fmt.Println("Found test bucket, removing bucket for testing.")
			DeleteBucket(svc, "example-bucket-for-test")

		}
	}

}

//GetSingleObject is a service that will retrive a single guest from the bucket
func GetSingleObject(svc *s3.S3, bucketName string, keyName string) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)

	svc = s3.New(sess)

	fmt.Println("Getting single Object: ", keyName)
	fmt.Println("from bucket: ", bucketName)

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	var responseBody Guest

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		fmt.Println("Error - Failed to read the response body. ", err)
		return
	}
	if err = json.Unmarshal(body, &responseBody); err == nil {
		fmt.Println(responseBody)
		fmt.Println(responseBody.Name)
	} else {
		fmt.Println("Error - Failed to read the response body. ", err)
		return
	}

	fmt.Println(result)

}

//Guest is guest struct that we want to insert into the S3 bucket
type Guest struct {
	Name      string
	Attending bool
	Cocktail  string
	Address   string
	Message   string
}

//addItemToBucket is a service that will take in a config and bucket name and add an item to it
func AddItemToBucket(svc *s3.S3, bucketName string) {
	fmt.Println("Welcome to AWS insert bucket test")

	g := Guest{"Solomon", true, "Phill's Cocktail", "100 Test Rd. Dothan, AL 36303", "Good Luck"}

	JSONGuest, err := json.Marshal(g)

	input := &s3.PutObjectInput{
		ACL:         aws.String("public-read"),
		Body:        aws.ReadSeekCloser(strings.NewReader(string(JSONGuest))),
		Bucket:      aws.String("example-bucket-for-test"),
		ContentType: aws.String("application/json"),
		Key:         aws.String(g.Name),
	}

	result, err := svc.PutObject(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)

}

//ListobjectsInBucket is a service that lists the objects in a given bucket
func ListobjectsInBucket(svc *s3.S3, bucketName string) {

	fmt.Println("Starting list of bucket objects")

	//we could list objects and then select them
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucketName),
		MaxKeys: aws.Int64(1000),
	}

	result, err := svc.ListObjectsV2(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)

}

//CreateBucket services used to create a bucket in AWS account
func CreateBucket(svc *s3.S3) {
	fmt.Println("Welcome to AWS create bucket test")

	input := &s3.CreateBucketInput{
		Bucket: aws.String("example-bucket-for-test"),
	}
	result, err := svc.CreateBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)

}

//DeleteBucket is a service that is used to remove buckets from the S3
func DeleteBucket(svc *s3.S3, bucketName string) {

	fmt.Println("Starting AWS Delete Method, removing bucket: ", bucketName)

	input := &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	}

	_, err := svc.DeleteBucket(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
}

//CreateAWSConfig services used to create a config and session in AWS account
func CreateAWSConfig() *s3.S3 {

	fmt.Println("Welcome to AWS CONFIG setup")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)

	if err != nil {
		fmt.Printf("Error creating the S# config.")
	}

	svc := s3.New(sess)

	return svc

}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
