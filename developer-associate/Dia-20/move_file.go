package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client *s3.Client
)

type Record struct {
	Bucket      string `json:"bucket"`
	Key         string `json:"key"`
	Destination string `json:"destination"`
}

type Output struct {
	Status      string `json:"status"`
	Destination string `json:"destination"`
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

	pathKey := strings.Split(record.Key, "/")

	destinationKey := fmt.Sprintf("%s/%s", record.Destination, pathKey[len(pathKey)-1])

	pathCopySource := fmt.Sprintf("%s/%s", record.Bucket, record.Key)

	s3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(record.Bucket),
		Key:        aws.String(destinationKey),
		CopySource: aws.String(pathCopySource),
	})

	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(record.Bucket),
		Key:    aws.String(record.Key),
	})

	return Output{
		Status:      "moved",
		Destination: destinationKey,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
