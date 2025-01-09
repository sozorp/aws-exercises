package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Transaction struct {
	TransactionID   string  `json:"transaction_id" gorm:"primaryKey;column:transaction_id;type:varchar(50)"`
	Amount          float64 `json:"amount" gorm:"not null;column:amount;type:decimal(10,2)"`
	TransactionType string  `json:"transaction_type" gorm:"not null;column:transaction_type;type:varchar(20)"`
	Timestamp       string  `json:"timestamp" gorm:"not null;column:timestamp"`
}

var (
	dbClient *gorm.DB
)

func init() {
	var err error

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	dbClient, err = gorm.Open(postgres.Open(dns), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
}

func CreateTransaction(transaction Transaction) error {

	if result := dbClient.Create(&transaction); result.Error != nil {
		return result.Error
	}

	return nil
}

func Handler(ctx context.Context, request events.KinesisEvent) error {

	for _, record := range request.Records {
		var transaction Transaction

		if err := json.Unmarshal(record.Kinesis.Data, &transaction); err != nil {
			log.Println(err)
			return fmt.Errorf("Error parsing transaction")
		}

		if err := CreateTransaction(transaction); err != nil {
			log.Println(err)
			return fmt.Errorf("Error saving transaction")
		}
	}

	log.Println("Transactions processed successfully")

	return nil
}

func main() {
	lambda.Start(Handler)
}
