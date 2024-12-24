package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, event json.RawMessage) (Output, error) {
	log.Println("Validating order:", event)

	order, err := ValidateOrder(event)

	if err != nil {
		log.Fatalln(err)

		return Output{}, err
	}

	return Output{Status: "Order Validated", OrderID: order.OrderID}, nil
}

func main() {
	lambda.Start(Handler)
}
