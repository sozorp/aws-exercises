package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
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

func Handler(ctx context.Context, event json.RawMessage) (Output, error) {
	var record Record

	if err := json.Unmarshal(event, &record); err != nil {
		return Output{}, fmt.Errorf("Error: decoding response %v", err)
	}

	log.Printf("Processing file: s3://%s/%s", record.Bucket, record.Key)

	return Output{
		Status: "processed",
		Bucket: record.Bucket,
		Key:    record.Key,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
