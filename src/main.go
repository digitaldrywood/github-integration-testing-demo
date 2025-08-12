package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	bucket   = flag.String("bucket", os.Getenv("S3_BUCKET"), "S3 bucket name")
	region   = flag.String("region", os.Getenv("AWS_REGION"), "AWS region")
	endpoint = flag.String("endpoint", os.Getenv("S3_ENDPOINT"), "S3 endpoint (for testing)")
)

func main() {
	flag.Parse()

	if *bucket == "" {
		log.Fatal("bucket is required (use -bucket flag or S3_BUCKET env var)")
	}

	// Create AWS session
	config := &aws.Config{}
	if *region != "" {
		config.Region = aws.String(*region)
	}
	if *endpoint != "" {
		config.Endpoint = aws.String(*endpoint)
		config.S3ForcePathStyle = aws.Bool(true)
	}

	sess, err := session.NewSession(config)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	// List objects in bucket
	fmt.Printf("Listing objects in bucket: %s\n", *bucket)
	result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(*bucket),
	})
	if err != nil {
		log.Fatalf("Failed to list objects: %v", err)
	}

	fmt.Printf("Found %d objects:\n", len(result.Contents))
	for _, item := range result.Contents {
		fmt.Printf("  %s (size: %d)\n", *item.Key, *item.Size)
	}
}

// UploadFile uploads a file to S3
func UploadFile(svc *s3.S3, bucket, key, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	fmt.Printf("Successfully uploaded %s to %s/%s\n", filename, bucket, key)
	return nil
}

// DownloadFile downloads a file from S3
func DownloadFile(svc *s3.S3, bucket, key, filename string) error {
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer result.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.ReadFrom(result.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Successfully downloaded %s/%s to %s\n", bucket, key, filename)
	return nil
}