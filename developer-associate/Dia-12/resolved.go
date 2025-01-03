package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Users struct {
	ID    uint   `json:"id" gorm:"primaryKey;column:id;autoIncrement"`
	Email string `json:"email" gorm:"not null;unique"`
	Name  string `json:"name" gorm:"not null"`
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

	dbClient.AutoMigrate(&Users{})
}

func HandlerGetAllUsers() (events.APIGatewayProxyResponse, error) {
	var users []Users

	if result := dbClient.Find(&users); result.Error != nil {
		log.Println("Error:", result.Error)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `"error":"There was an error getting the records"`,
		}, nil
	}

	response, err := json.Marshal(users)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Failed to encode response"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil
}

func HandlerCreateUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var user Users

	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		log.Println(err)

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"error": "Invalid request"}`),
		}, nil
	}

	if result := dbClient.Create(&user); result.Error != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Could not create user"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       fmt.Sprintf(`{"message": "User created", "id": %d}`, user.ID),
	}, nil
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println(request.HTTPMethod)
	switch request.HTTPMethod {
	case http.MethodGet:
		return HandlerGetAllUsers()
	case http.MethodPost:
		return HandlerCreateUser(request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       `{"error": "Method not allowed"}`,
		}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
