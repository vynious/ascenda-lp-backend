package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/vynious/ascenda-lp-backend/db"
	"github.com/vynious/ascenda-lp-backend/types"
	"gorm.io/gorm"
)

var (
	DBService *db.DBService
	RDSClient *rds.Client
	err       error
)

func init() {
	DBService, err = db.SpawnDBService()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if request.RequestContext.HTTP.Method == "OPTIONS" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type",
			},
			Body: "{}",
		}, nil
	}

	var userRequestBody types.UpdateUserRequestBody

	if err := json.Unmarshal([]byte(request.Body), &userRequestBody); err != nil {
		log.Printf("JSON unmarshal error: %s", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "PUT",
			},
			Body: "Invalid request format",
		}, nil
	}

	updatedUser, err := db.UpdateUserWithUpdateUserRequestBody(ctx, DBService, userRequestBody)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 404,
				Headers: map[string]string{
					"Access-Control-Allow-Headers": "Content-Type",
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "PUT",
				},
				Body: "User not found",
			}, nil
		} else {
			log.Printf("Database error: %s", err)
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 500,
				Headers: map[string]string{
					"Access-Control-Allow-Headers": "Content-Type",
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "PUT",
				},
				Body: "Internal server error",
			}, nil
		}
	}

	responseBody, err := json.Marshal(updatedUser)
	if err != nil {
		log.Printf("JSON marshal error: %s", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "PUT",
			},
			Body: "Error marshaling user into JSON",
		}, nil
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "PUT",
		},
		Body: string(responseBody),
	}, nil
}

func main() {
	// we are simulating a lambda behind an ApiGatewayV2
	lambda.Start(handler)
	defer DBService.CloseConn()
}
