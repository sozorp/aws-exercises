package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

const (
	tableName = "EventLogs"
)

var (
	dynamoClient *dynamodb.Client
)

type Data struct {
	Sensor string  `dynamodbav:"sensor" json:"sensor"`
	Value  float64 `dynamodbav:"value" json:"value"`
}

type Event struct {
	EventID   string `dynamodbav:"eventId" json:"eventId"`
	Timestamp string `dynamodbav:"timestamp" json:"timestamp"`
	Data      Data   `dynamodbav:"data" json:"data"`
}

func GetEvents() ([]Event, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynamoClient.Scan(context.TODO(), input)

	if err != nil {
		return nil, fmt.Errorf("error scanning table %s", tableName)
	}

	var events []Event

	err = attributevalue.UnmarshalListOfMaps(result.Items, &events)

	if err != nil {
		return nil, fmt.Errorf("error deserializing data %v", err)
	}

	return events, nil
}

func HomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to the Gin APP on Elastic Beanstalk!")
}

func GetEventsHandler(c *gin.Context) {
	items, err := GetEvents()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))

	if err != nil {
		log.Panicf("unable to load SDK config, %v", err)
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)

	router := gin.Default()

	router.GET("/", HomeHandler)
	router.GET("/events", GetEventsHandler)

	router.Run(":5000")
}
