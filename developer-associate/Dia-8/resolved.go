package main

import (
	"context"
	"encoding/json"
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

var (
	dynamoClient *dynamodb.Client
)

type Data struct {
	Sensor string  `dynamodbav:"sensor" json:"sensor"`
	Value  float64 `dynamodbav:"value" json:"value"`
}

type Record struct {
	EventID   string `dynamodbav:"eventId" json:"eventId"`
	Timestamp string `dynamodbav:"timestamp" json:"timestamp"`
	Data      Data   `dynamodbav:"data" json:"data"`
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
}

func uploadRecordToDynamo(ctx context.Context, tableName string, item Record) error {
	av, err1 := attributevalue.MarshalMap(item)

	if err1 != nil {
		return err1
	}

	input := dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	if _, err := dynamoClient.PutItem(ctx, &input); err != nil {
		return err
	}

	return nil
}

func Handler(ctx context.Context, event events.KinesisEvent) error {
	tableName := os.Getenv("TABLE_NAME")

	for _, item := range event.Records {

		kinesisData := item.Kinesis.Data

		var payload Data

		err2 := json.Unmarshal(kinesisData, &payload)

		if err2 != nil {
			log.Fatalln(err2)
			return err2
		}

		err := uploadRecordToDynamo(ctx, tableName, Record{
			EventID:   uuid.New().String(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Data:      payload,
		})

		if err != nil {
			log.Fatalln(err)
			return err
		}

		log.Printf("Processed event %+v", payload)
	}

	log.Println("Events processed sucessfully")
	return nil
}

func main() {
	lambda.Start(Handler)
}
