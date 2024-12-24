package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Body struct {
	Type      string `json:"Type"`
	MessageId string `json:"MessageId"`
	Message   string `json:"Message"`
}

func Handler(ctx context.Context, event events.SQSEvent) error {

	for _, message := range event.Records {

		body := Body{}

		if err := json.Unmarshal([]byte(message.Body), &body); err != nil {
			log.Fatalln("Error al deserializar el mensaje:", err)
			return err
		}

		log.Printf("Mensaje recibido de tipo: %s \n", body.Type)
		log.Printf("Procesando mensaje: %s\n", body.Message)
		log.Printf("El mensaje fue procesado correctamente con id: %s", body.MessageId)
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
