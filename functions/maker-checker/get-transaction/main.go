package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vynious/ascenda-lp-backend/db"
	"github.com/vynious/ascenda-lp-backend/types"
	"log"
)

var (
	DBService    *db.DBService
	responseBody types.TransactionResponseBody

	requestBody types.GetTransactionRequestBody
	err         error
)

func init() {

	// Initialise global variable DBService tied to Aurora
	DBService, err = db.SpawnDBService()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func GetTransactionHandler(ctx context.Context, req *events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	defer DBService.CloseConn()

	if err := json.Unmarshal([]byte(req.Body), &requestBody); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 404,
			Body:       "Bad Request",
		}, nil
	}

	txn, err := DBService.GetTransaction(ctx, requestBody.TransactionId)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       "",
		}, nil
	}

	responseBody.Txn = *txn

	bod, err := json.Marshal(responseBody)

	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 201,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(bod),
	}, nil
}

func main() {
	lambda.Start(GetTransactionHandler)
}
