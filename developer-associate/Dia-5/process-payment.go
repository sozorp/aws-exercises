package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"atomicgo.dev/random"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, event json.RawMessage) (Output, error) {
	log.Println("Processing payment for order:", event)

	order, err := ValidateOrder(event)

	if err != nil {
		log.Fatalln(err)

		return Output{}, err
	}

	if random.Bool() {
		return Output{}, fmt.Errorf("Payment processing failed")
	}

	return Output{Status: "Payment Successful", OrderID: order.OrderID}, nil
}

func main() {
	lambda.Start(Handler)
}
