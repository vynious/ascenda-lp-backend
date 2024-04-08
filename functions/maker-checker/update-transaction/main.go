package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vynious/ascenda-lp-backend/db"
	makerchecker "github.com/vynious/ascenda-lp-backend/types"
	"github.com/vynious/ascenda-lp-backend/util"
)

var (
	DBService    *db.DBService
	requestBody  makerchecker.UpdateTransactionRequestBody
	responseBody makerchecker.TransactionResponseBody
	err          error
	DB           *db.DB
)

func init() {
	DBService, err = db.SpawnDBService()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func LambdaHandler(ctx context.Context, req *events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	// Checking if userid and userlocation exists for logging purposes
	userId, err := util.GetCustomAttributeWithCognito("custom:userId", req.Headers["Authorization"])
	if err != nil {
		ctx = context.WithValue(ctx, "userId", userId)
	}
	userLocation, ok := req.Headers["CloudFront-Viewer-Country"]
	if ok {
		ctx = context.WithValue(ctx, "userLocation", userLocation)
	}
	bank, err := util.GetCustomAttributeWithCognito("custom:bank", req.Headers["Authorization"])
	if err != nil {
		ctx = context.WithValue(ctx, "bank", bank)
	}
	/*
		check role/user of requested
	*/

	DB = DBService.GetBanksDB(req.Headers["Authorization"])

	if err := json.Unmarshal([]byte(req.Body), &requestBody); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Headers: map[string]string{
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "PATCH",
			},
			Body: "Bad Request",
		}, nil
	}

	updatedTxn, err := DB.UpdateTransaction(ctx, requestBody.TransactionId, requestBody.CheckerId, requestBody.Approval)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "PATCH",
			},
			Body: "",
		}, nil
	}
	responseBody.Txn = *updatedTxn

	bod, err := json.Marshal(responseBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 201,
			Headers: map[string]string{
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "PATCH",
			},
			Body: err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers: map[string]string{
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "PATCH",
		},
		Body: string(bod),
	}, nil
}

func main() {
	lambda.Start(LambdaHandler)
	defer DBService.CloseConnections()
}
