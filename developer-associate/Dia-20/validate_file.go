package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client *s3.Client
)

type Record struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

type Output struct {
	Status string `json:"status"`
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Panicf("unable to load SDK config, %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
}

func Handler(ctx context.Context, event json.RawMessage) (Output, error) {

	var record Record

	if err := json.Unmarshal(event, &record); err != nil {
		return Output{}, fmt.Errorf("Error: decoding response %v", err)
	}

	response, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(record.Bucket),
		Key:    aws.String(record.Key),
	})

	if err != nil {
		return Output{}, err
	}

	if *response.ContentLength > 1*1024*1024 {
		return Output{}, fmt.Errorf("File size exceeds the allowed limit")
	}

	return Output{
		Status: "valid",
		Bucket: record.Bucket,
		Key:    record.Key,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
