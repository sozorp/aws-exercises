package main

import (
	"context"
	"fmt"
	"log"

	"atomicgo.dev/random"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, event events.S3Event) error {

	for _, record := range event.Records {

		objectKey := record.S3.Object.Key

		log.Println("Processing the file with the key: ", objectKey)

		if random.Bool() {
			return fmt.Errorf("The file does not meet expectations")
		}

		log.Printf("The %s file was processed successfully", objectKey)
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
