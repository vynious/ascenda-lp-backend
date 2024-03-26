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
	"github.com/joho/godotenv"
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
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env")
	}

	DBService, err = db.SpawnDBService()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var roleRequestBody types.CreateRoleRequestBody

	if err := json.Unmarshal([]byte(request.Body), &roleRequestBody); err != nil {
		log.Printf("JSON unmarshal error: %s", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       "Invalid request format",
		}, nil
	}

	_, err := db.RetrieveRoleWithRoleName(ctx, DBService, roleRequestBody.RoleName)
	if err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			roleName, err := db.CreateRoleWithCreateRoleRequestBody(ctx, DBService, roleRequestBody)
			if err != nil {
				log.Printf("Database error: %s", err)
				return events.APIGatewayV2HTTPResponse{
					StatusCode: 500,
					Body:       "Internal server error",
				}, nil
			}

			responseBody := fmt.Sprintf("{\"role_name\": \"%s\"}", roleName)
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 200,
				Body:       responseBody,
			}, nil
		} else {
			log.Printf("Database error: %s", err)
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 500,
				Body:       "Internal server error",
			}, nil
		}
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 409,
		Body:       "Role already exist. Please use a different name",
	}, nil
}

func main() {
	// we are simulating a lambda behind an ApiGatewayV2
	lambda.Start(handler)
	defer DBService.CloseConn()
}
