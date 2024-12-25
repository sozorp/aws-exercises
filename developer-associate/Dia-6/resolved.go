package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Claims struct {
	Email    string `json:"email"`
	Sub      string `json:"sub"`
	Nickname string `json:"nickname"`
	Name     string `json:"name"`
}

func InterfaceToStruct(input interface{}, output interface{}) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("error al convertir a JSON: %w", err)
	}

	err = json.Unmarshal(jsonData, output)
	if err != nil {
		return fmt.Errorf("error al parsear el JSON: %w", err)
	}

	return nil
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	authContext, err := request.RequestContext.Authorizer["claims"]

	if !err {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       `{"status": "error", "message": "Unauthorized"}`,
		}, nil
	}

	var claims Claims

	err1 := InterfaceToStruct(authContext, &claims)

	if err1 != nil {
		log.Printf("Error marshaling response: %v", err1)

		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"status": "error", "message": "Internal Server Error"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Hello, %s! (%s) Your user ID is %s, and email is %s", claims.Name, claims.Nickname, claims.Sub, claims.Email),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
