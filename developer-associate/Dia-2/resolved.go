package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Log struct {
	FileName  string `json:"file_name"`
	FileSize  int64  `json:"file_size"`
	Timestamp string `json:"timestamp"`
}

var (
	s3Client *s3.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
}

func uploadLogToS3(ctx context.Context, body []byte, bucket, key string) error {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(body),
		ContentType: aws.String("application/json"),
	}

	if _, err := s3Client.PutObject(ctx, input); err != nil {
		log.Printf("Failed to save log to S3: %v", err)
		return err
	}

	return nil
}

func Handler(ctx context.Context, event events.S3Event) error {
	data := event.Records[0].S3

	bucketName := data.Bucket.Name
	objectKey := data.Object.Key
	objectKeySlice := strings.Split(objectKey, "/")
	fileName := objectKeySlice[len(objectKeySlice)-1]
	log_key := fmt.Sprintf("logs/%s.json", fileName)

	log := Log{
		FileName:  objectKey,
		FileSize:  data.Object.Size,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	body, err := json.Marshal(log)

	if err != nil {
		return err
	}

	if err = uploadLogToS3(ctx, body, bucketName, log_key); err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
