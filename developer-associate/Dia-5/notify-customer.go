package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, event json.RawMessage) (Output, error) {
	log.Println("Notifying customer for order:", event)

	order, err := ValidateOrder(event)

	if err != nil {
		log.Fatalln(err)

		return Output{}, err
	}

	return Output{Status: "Notification Sent", OrderID: order.OrderID}, nil
}

func main() {
	lambda.Start(Handler)
}
