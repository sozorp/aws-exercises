package main

import (
	"context"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func fibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}

	var n2, n1 = 0, 1

	for i := 2; i <= n; i++ {
		n2, n1 = n1, n1+n2
	}

	return n1
}

func Handler(ctx context.Context, event events.SQSEvent) error {

	for _, message := range event.Records {
		messageBody := message.Body

		log.Println("Processing message: ", messageBody)

		number, err := strconv.Atoi(messageBody)

		if err != nil {
			log.Fatalln("Error parsing the number, it is invalid")

			return err
		}

		// heavy work simulation
		result := fibonacciIterative(number)

		log.Println("Result: ", result)
		log.Println("Task completed successfully!")
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
