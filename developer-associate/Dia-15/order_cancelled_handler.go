package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var (
	snsClient *sns.Client
)

type Order struct {
	OrderId string `json:"orderId"`
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Panicf("unable to load SDK config, %v", err)
	}

	snsClient = sns.NewFromConfig(cfg)
}

func SendSNS(ctx context.Context, message, subject, topicArn string) error {
	input := sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Subject:  aws.String(subject),
		Message:  aws.String(message),
	}

	_, err := snsClient.Publish(ctx, &input)

	if err != nil {
		return fmt.Errorf("There was an error sending the message, %v", err)
	}

	return nil
}

func Handler(ctx context.Context, request events.EventBridgeEvent) error {
	topicArn := os.Getenv("SNS_TOPIC_ARN")

	var order Order

	if err := json.Unmarshal(request.Detail, &order); err != nil {
		log.Println("Error: ", err)
		return err
	}

	message := fmt.Sprintf("Order %s has been cancelled.", order.OrderId)

	err := SendSNS(ctx, message, "Order Cancelled", topicArn)

	return err
}

func main() {
	lambda.Start(Handler)
}
