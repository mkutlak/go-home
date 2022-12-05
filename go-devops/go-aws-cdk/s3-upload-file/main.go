package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const bucketName = "aws-mkutlak-test-demo-bucket-go"
const region = "eu-west-1"

func main() {
	var (
		s3Client *s3.Client
		err      error
	)
	ctx := context.Background()

	if s3Client, err = initS3Client(ctx); err != nil {
		fmt.Printf("initS3Client error: %v\n", err)
		os.Exit(1)
	}

	if err = createS3Bucket(ctx, s3Client); err != nil {
		fmt.Printf("createS3Bucket error: %v\n", err)
		os.Exit(1)
	}

	if err = uploadToS3Bucket(ctx, s3Client); err != nil {
		fmt.Printf("uploadToS3Bucket error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Upload completed!")
}

func initS3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))

	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	return s3.NewFromConfig(cfg), nil
}

func createS3Bucket(ctx context.Context, s3Client *s3.Client) error {
	allBuckets, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})

	if err != nil {
		return fmt.Errorf("unable to list S3 buckets: %v", err)
	}

	var found bool = false

	for _, bucket := range allBuckets.Buckets {
		if *bucket.Name == bucketName {
			found = true
		}
	}
	if !found {
		_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: region,
			},
		})

		if err != nil {
			return fmt.Errorf("unable create S3 bucket: %v", err)
		}
	}

	return nil
}

func uploadToS3Bucket(ctx context.Context, s3client *s3.Client) error {

	uploader := manager.NewUploader(s3client)
	textFile, err := ioutil.ReadFile("text.file.log")

	if err != nil {
		return fmt.Errorf("unable read from file: %v", err)
	}
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("text-file.txt"),
		Body:   bytes.NewReader(textFile),
	})

	if err != nil {
		return fmt.Errorf("unable upload file to S3 bucket: %v", err)
	}

	return nil
}
