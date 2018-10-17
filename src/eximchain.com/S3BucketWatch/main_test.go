package main

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// TestS3Handler is used to test that the logic for the S3Handler is coded correctly.
// The notification bucket and the key affected needs to be updated and tested manually.
func TestS3Handler(t *testing.T) {
	var (
		ctx     context.Context
		s3Event events.S3Event
	)

	s3Event.Records = make([]events.S3EventRecord, 1)
	s3Event.Records[0].EventSource = "TestS3Handler"
	s3Event.Records[0].EventTime = time.Now()
	s3Event.Records[0].S3.Bucket.Name = "terraform-20181001033020781800000004" // the notification bucket
	s3Event.Records[0].S3.Object.Key = "BucketKey"                             // the key affected

	S3Handler(ctx, s3Event)
}

// BucketUpload uploads the specified data to the BucketName, Key in the given region
// Example usage as follows:
// region := ""
// BucketName := os.Args[1]
// BucketKey := os.Args[2]
// if len(os.Args) > 3 {
// 	region = os.Args[3]
// }
// BucketUpload(bytes.NewReader([]byte("Hello world")), BucketName, BucketKey, region)
func BucketUpload(reader *bytes.Reader, BucketName, BucketKey, region string) {
	var sess *session.Session

	if sess == nil {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}
	var s3svc *s3.S3
	if region == "" {
		s3svc = s3.New(sess)
	} else {
		s3svc = s3.New(sess, aws.NewConfig().WithMaxRetries(3).WithRegion(region))
	}

	size := int64(reader.Size())
	if _, err := s3svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(BucketName),
		Key:           aws.String(BucketKey),
		ACL:           aws.String("bucket-owner-full-control"),
		Body:          reader,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String("application/octet-stream"),
	}); err == nil {
		fmt.Println("Uploaded successfully")
	} else {
		fmt.Printf("Failed to upload, error: %v\n", err)
	}
}
