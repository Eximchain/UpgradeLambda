package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Handler to handle the SNS topic for a S3 event
func S3Handler(ctx context.Context, s3Event events.S3Event) {
	var sess *session.Session
	for _, record := range s3Event.Records {
		s3entity := record.S3
		BucketName := s3entity.Bucket.Name
		BucketKey := s3entity.Object.Key
		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, BucketName, BucketKey)

		if sess == nil {
			sess = session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))
		}

		// Create ec2 svc
		if ec2svc := ec2.New(sess); ec2svc != nil {

		}
		// Create s3svc
		if s3svc := s3.New(sess); s3svc != nil {

			downloader := s3manager.NewDownloader(sess)

			getS3BucketObj := func(downloader *s3manager.Downloader, bucketName, bucketKey string) (result []byte, err error) {
				var n int64
				buffer := []byte{}
				buf := aws.NewWriteAtBuffer(buffer)
				fmt.Printf("Downloading %s %s\n", bucketName, bucketKey)
				if n, err = downloader.Download(buf, &s3.GetObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(bucketKey),
				}); err == nil {
					result = buf.Bytes()
					fmt.Printf("Bytes read: %v for Bucket: %s, Key: %s\n", n, bucketName, bucketKey)
				} else {
					fmt.Printf("Bucket: %s, Item: %s, err: %v\n", bucketName, bucketKey, err)
				}
				return
			}

			if bucketKeyBytes, err := getS3BucketObj(downloader, BucketName, BucketKey); err == nil {
				bucketKeyContent := string(bucketKeyBytes)
				fmt.Printf("Content of Bucket: %s, Key: %s: %s", BucketName, BucketKey, bucketKeyContent)

			} else {
				fmt.Printf("Failed to retrieve Bucket: %s, Key: %s", BucketName, BucketKey)
			}

			if _, err := s3svc.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(BucketName),
				Key:    aws.String(BucketKey),
			}); err == nil {
				fmt.Printf("Bucket: %s, Key: %s deleted", BucketName, BucketKey)
			} else {
				fmt.Printf("Failed to delete Bucket: %s, Key: %s", BucketName, BucketKey)
			}
		}

	}
}

func main() {
	lambda.Start(S3Handler)
}
