// +build integration

package tests

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	runS3Tests = flag.Bool("s3", false, "Run S3 integration tests")
	runDBTests = flag.Bool("database", false, "Run database integration tests")
	runAPITests = flag.Bool("api", false, "Run API integration tests")
	
	s3Bucket = flag.String("bucket", os.Getenv("TEST_S3_BUCKET"), "S3 bucket for testing")
	s3Region = flag.String("region", os.Getenv("AWS_REGION"), "AWS region")
	s3Endpoint = flag.String("endpoint", os.Getenv("S3_ENDPOINT"), "S3 endpoint")
)

func TestS3Upload(t *testing.T) {
	if !*runS3Tests {
		t.Skip("S3 tests not enabled (use -s3 flag)")
	}

	if *s3Bucket == "" {
		t.Fatal("TEST_S3_BUCKET environment variable or -bucket flag required")
	}

	// Setup
	sess, svc := setupS3(t)
	testKey := fmt.Sprintf("test-upload-%d.txt", time.Now().Unix())
	testData := []byte("Hello, S3 Integration Test!")

	// Upload
	t.Logf("Uploading to s3://%s/%s", *s3Bucket, testKey)
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(*s3Bucket),
		Key:    aws.String(testKey),
		Body:   bytes.NewReader(testData),
	})
	if err != nil {
		t.Fatalf("Failed to upload object: %v", err)
	}

	// Verify upload
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(*s3Bucket),
		Key:    aws.String(testKey),
	})
	if err != nil {
		t.Fatalf("Failed to retrieve uploaded object: %v", err)
	}
	defer result.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(result.Body)
	if !bytes.Equal(buf.Bytes(), testData) {
		t.Errorf("Retrieved data doesn't match. Got %s, want %s", buf.String(), string(testData))
	}

	// Cleanup
	t.Cleanup(func() {
		_, err := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(*s3Bucket),
			Key:    aws.String(testKey),
		})
		if err != nil {
			t.Logf("Warning: Failed to cleanup test object: %v", err)
		}
	})

	t.Logf("✅ S3 upload test passed")
}

func TestS3ListObjects(t *testing.T) {
	if !*runS3Tests {
		t.Skip("S3 tests not enabled (use -s3 flag)")
	}

	if *s3Bucket == "" {
		t.Fatal("TEST_S3_BUCKET environment variable or -bucket flag required")
	}

	_, svc := setupS3(t)

	// List objects
	t.Logf("Listing objects in s3://%s", *s3Bucket)
	result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:  aws.String(*s3Bucket),
		MaxKeys: aws.Int64(10),
	})
	if err != nil {
		t.Fatalf("Failed to list objects: %v", err)
	}

	t.Logf("Found %d objects in bucket", len(result.Contents))
	for _, obj := range result.Contents {
		t.Logf("  - %s (size: %d bytes)", *obj.Key, *obj.Size)
	}

	t.Logf("✅ S3 list objects test passed")
}

func TestS3BucketVersioning(t *testing.T) {
	if !*runS3Tests {
		t.Skip("S3 tests not enabled (use -s3 flag)")
	}

	if *s3Bucket == "" {
		t.Fatal("TEST_S3_BUCKET environment variable or -bucket flag required")
	}

	_, svc := setupS3(t)

	// Check versioning status
	t.Logf("Checking versioning status for s3://%s", *s3Bucket)
	versioning, err := svc.GetBucketVersioning(&s3.GetBucketVersioningInput{
		Bucket: aws.String(*s3Bucket),
	})
	if err != nil {
		t.Fatalf("Failed to get bucket versioning: %v", err)
	}

	if versioning.Status != nil {
		t.Logf("Bucket versioning status: %s", *versioning.Status)
	} else {
		t.Log("Bucket versioning is not enabled")
	}

	t.Logf("✅ S3 versioning check passed")
}

func TestDatabaseConnection(t *testing.T) {
	if !*runDBTests {
		t.Skip("Database tests not enabled (use -database flag)")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		t.Skip("DB_HOST not set, skipping database tests")
	}

	// This is a placeholder for database testing
	// In a real scenario, you would:
	// 1. Connect to the database
	// 2. Run migrations
	// 3. Perform CRUD operations
	// 4. Verify data integrity

	t.Logf("Would connect to database at %s", dbHost)
	t.Logf("✅ Database connection test passed (mock)")
}

func TestExternalAPI(t *testing.T) {
	if !*runAPITests {
		t.Skip("API tests not enabled (use -api flag)")
	}

	apiEndpoint := os.Getenv("API_ENDPOINT")
	if apiEndpoint == "" {
		t.Skip("API_ENDPOINT not set, skipping API tests")
	}

	// This is a placeholder for API testing
	// In a real scenario, you would:
	// 1. Make HTTP requests to the API
	// 2. Verify response codes
	// 3. Validate response data
	// 4. Test error handling

	t.Logf("Would test API at %s", apiEndpoint)
	t.Logf("✅ API test passed (mock)")
}

// Helper function to setup S3 client
func setupS3(t *testing.T) (*session.Session, *s3.S3) {
	config := &aws.Config{}
	
	if *s3Region != "" {
		config.Region = aws.String(*s3Region)
	}
	
	if *s3Endpoint != "" {
		config.Endpoint = aws.String(*s3Endpoint)
		config.S3ForcePathStyle = aws.Bool(true)
		t.Logf("Using custom S3 endpoint: %s", *s3Endpoint)
	}

	sess, err := session.NewSession(config)
	if err != nil {
		t.Fatalf("Failed to create AWS session: %v", err)
	}

	return sess, s3.New(sess)
}