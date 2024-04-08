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
	"github.com/vynious/ascenda-lp-backend/util"
	"gorm.io/gorm"
)

var (
	DBService *db.DBService
	RDSClient *rds.Client
	err       error
	DB        *db.DB
)

func init() {
	DBService, err = db.SpawnDBService()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("GetUser")
	// Checking if userid and userlocation exists for logging purposes
	userId, err := util.GetCustomAttributeWithCognito("custom:userID", request.Headers["Authorization"])
	if err == nil {
		log.Printf("GetCustomAttribute custom:userID %s", userId)
		ctx = context.WithValue(ctx, "userId", userId)
	}
	userLocation, ok := request.Headers["CloudFront-Viewer-Country"]
	if ok {
		log.Printf("Get Attribute CloudFront-Viewer-Country")
		ctx = context.WithValue(ctx, "userLocation", userLocation)
	}
	bank, err := util.GetCustomAttributeWithCognito("custom:bank", request.Headers["Authorization"])
	if err == nil {
		log.Printf("GetCustomAttribute custom:bank %s", bank)
		ctx = context.WithValue(ctx, "bank", bank)
	}

	DB = DBService.GetBanksDB(request.Headers["Authorization"])

	users, err := db.RetrieveAllUsers(ctx, DB)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Headers: map[string]string{
					"Access-Control-Allow-Headers": "Content-Type",
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "GET",
				},
				Body: "User(s) not found",
			}, nil
		} else {
			log.Printf("Database error: %s", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers: map[string]string{
					"Access-Control-Allow-Headers": "Content-Type",
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "GET",
				},
				Body: `{"message": "Internal server error"}`,
			}, nil
		}
	}

	responseBody, err := json.Marshal(users)
	if err != nil {
		log.Printf("JSON marshal error: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET",
			},
			Body: `{"message": "Error marshaling users into JSON"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET",
		},
		Body: string(responseBody),
	}, nil
}

func main() {
	// we are simulating a lambda behind an ApiGatewayV2
	lambda.Start(handler)
	defer DBService.CloseConnections()
}
