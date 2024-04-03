package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	var userRequestBody types.CreateUserRequestBody

	if err := json.Unmarshal([]byte(request.Body), &userRequestBody); err != nil {
		log.Printf("JSON unmarshal error: %s", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "POST",
			},
			Body: "Invalid request format",
		}, nil
	}

	_, err := db.RetrieveUserWithEmail(ctx, DBService, userRequestBody.Email)
	if err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user, err := db.CreateUserWithCreateUserRequestBody(ctx, DBService, userRequestBody)
			if err != nil {
				log.Printf("Database error: %s", err)
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return events.APIGatewayV2HTTPResponse{
						StatusCode: 404,
						Headers: map[string]string{
							"Access-Control-Allow-Headers": "Content-Type",
							"Access-Control-Allow-Origin":  "*",
							"Access-Control-Allow-Methods": "POST",
						},
						Body: "Role not found. Please create a role first (if user have a role)",
					}, nil
				}
				return events.APIGatewayV2HTTPResponse{
					StatusCode: 500,
					Headers: map[string]string{
						"Access-Control-Allow-Headers": "Content-Type",
						"Access-Control-Allow-Origin":  "*",
						"Access-Control-Allow-Methods": "POST",
					},
					Body: "Internal server error",
				}, nil
			}

			responseBody := fmt.Sprintf("{\"email\": \"%s\", \"id\": \"%s\"}", user.Email, user.Id)
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"Access-Control-Allow-Headers": "Content-Type",
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "POST",
				},
				Body: responseBody,
			}, nil
		} else {
			log.Printf("Database error: %s", err)
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 500,
				Headers: map[string]string{
					"Access-Control-Allow-Headers": "Content-Type",
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "POST",
				},
				Body: "Internal server error",
			}, nil
		}
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 409,
		Headers: map[string]string{
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		},
		Body: "User already exist. Please use a different email",
	}, nil
}

func main() {
	lambda.Start(handler)
	defer DBService.CloseConn()
}
