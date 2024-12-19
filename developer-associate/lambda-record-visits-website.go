package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

type VisitRecord struct {
	VisitID   string `dynamodbav:"visit_id" json:"visit_id"`
	Timestamp string `dynamodbav:"timestamp" json:"timestamp"`
	UserAgent string `dynamodbav:"user_agent" json:"user_agent"`
}

var (
	dynamoClient *dynamodb.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
}

func uploadVisitRecordToDynamoDB(ctx context.Context, tableName string, item VisitRecord) error {

	av, err := attributevalue.MarshalMap(item)

	if err != nil {
		log.Fatalf("Got error marshalling new visit record item: %s", err)
	}

	if _, ok := av["visit_id"]; !ok {
		return fmt.Errorf("visit_id no está presente en el ítem: %+v", av)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}

	_, err = dynamoClient.PutItem(ctx, input)

	if err != nil {
		log.Printf("Failed to save visit record to DynamoDB: %v", err)
		return err
	}

	return nil
}

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var visit VisitRecord

	userAgent, ok := request.Headers["user-agent"]

	if !ok {
		userAgent = "Unknown"
	}

	visit.VisitID = uuid.New().String()
	visit.Timestamp = time.Now().UTC().String()
	visit.UserAgent = userAgent

	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		log.Printf("TABLE_NAME environment variable is not set")
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("missing required environment variable TABLE_NAME"),
		}, fmt.Errorf("missing required environment variable TABLE_NAME")
	}

	if err := uploadVisitRecordToDynamoDB(ctx, tableName, visit); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}

	log.Printf("Successfully save visit record %s in DynamoDB table: %s", visit.VisitID, tableName)

	responseBody, _ := json.Marshal(visit)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
