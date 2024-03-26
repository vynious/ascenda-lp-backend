package main

import (
	"context"
	"encoding/json"
	"errors"
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
	var roleRequestBody types.UpdateRoleRequestBody

	if err := json.Unmarshal([]byte(request.Body), &roleRequestBody); err != nil {
		log.Printf("JSON unmarshal error: %s", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       "Invalid request format",
		}, nil
	}

	role, err := db.RetrieveRoleWithRoleName(ctx, DBService, roleRequestBody.RoleName)
	if err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 404,
				Body:       "Role not found",
			}, nil
		} else {
			log.Printf("Database error: %s", err)
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 500,
				Body:       "Internal server error",
			}, nil
		}
	}

	db.UpdateRole(ctx, DBService, roleRequestBody)

	responseBody, err := json.Marshal(role)
	if err != nil {
		log.Printf("JSON marshal error: %s", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       "Error marshaling role into JSON",
		}, nil
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}
func main() {
	// we are simulating a lambda behind an ApiGatewayV2
	lambda.Start(handler)
	defer DBService.CloseConn()
}
