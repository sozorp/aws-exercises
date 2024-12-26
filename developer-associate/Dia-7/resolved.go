package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var (
	s3Client  *s3.Client
	snsClient *sns.Client
)

type Response struct {
	Status  uint   `json:"status"`
	Message string `json:"message"`
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	snsClient = sns.NewFromConfig(cfg)
}

func checkENVs() error {
	s3BucketName := os.Getenv("S3_BUCKET_NAME")

	if s3BucketName == "" {
		return fmt.Errorf("missing required environment variable S3_BUCKET_NAME")
	}

	snsArn := os.Getenv("SNS_ARN")

	if snsArn == "" {
		return fmt.Errorf("missing required environment variable SNS_ARN")
	}

	return nil
}

func uploadToS3(ctx context.Context, body []byte, bucket, key string) error {
	input := s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(body),
		ContentType: aws.String("text/plain"),
	}

	_, err := s3Client.PutObject(ctx, &input)

	if err != nil {
		return err
	}

	return nil
}

func publishToSNS(ctx context.Context, arn, message, subject string) error {

	input := sns.PublishInput{
		TopicArn: aws.String(arn),
		Message:  aws.String(message),
		Subject:  aws.String(subject),
	}

	_, err := snsClient.Publish(ctx, &input)

	if err != nil {
		return err
	}

	return nil
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	err := checkENVs()

	if err != nil {
		log.Fatal(err)

		body, _ := json.Marshal(Response{
			Status:  http.StatusInternalServerError,
			Message: "There was a problem, please try again later",
		})

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(body),
		}, err
	}

	userId, ok1 := event.QueryStringParameters["userId"]

	if !ok1 {
		body, _ := json.Marshal(Response{
			Status:  http.StatusBadRequest,
			Message: "Missing userId",
		})

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	reportType, ok2 := event.QueryStringParameters["reportType"]

	if !ok2 {
		body, _ := json.Marshal(Response{
			Status:  http.StatusBadRequest,
			Message: "Missing reportType",
		})

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	bucketName := os.Getenv("S3_BUCKET_NAME")

	fileContent := fmt.Sprintf("Report for %s\nType: %s\nGenerated at: %s", userId, reportType, time.Now().UTC())
	fileName := fmt.Sprintf("%s-%s-%d.txt", userId, reportType, time.Now().UTC().Unix())

	err = uploadToS3(ctx, []byte(fileContent), bucketName, fileName)

	if err != nil {
		log.Fatal(err)

		body, _ := json.Marshal(Response{
			Status:  http.StatusInternalServerError,
			Message: "There was an error creating the report",
		})

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(body),
		}, err
	}

	fileUrl := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, fileName)

	messageSNS := fmt.Sprintf("Hello %s, your report is ready: %s", userId, fileUrl)

	err = publishToSNS(ctx, os.Getenv("SNS_ARN"), messageSNS, "Your report is ready")

	if err != nil {
		log.Fatal(err)

		body, _ := json.Marshal(Response{
			Status:  http.StatusInternalServerError,
			Message: "There was an error creating the report",
		})

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(body),
		}, err
	}

	body, _ := json.Marshal(Response{
		Status:  http.StatusOK,
		Message: "The report was created correctly",
	})

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
