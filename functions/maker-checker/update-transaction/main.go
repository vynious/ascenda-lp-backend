package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vynious/ascenda-lp-backend/db"
	makerchecker "github.com/vynious/ascenda-lp-backend/types"
)

var (
	DBService    *db.DBService
	requestBody  makerchecker.UpdateTransactionRequestBody
	responseBody makerchecker.TransactionResponseBody
	err          error
)

func init() {

	// Initialise global variable DBService tied to Aurora
	DBService, err = db.SpawnDBService()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func LambdaHandler(ctx context.Context, req *events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	/*
		check role/user of requested
	*/

	if err := json.Unmarshal([]byte(req.Body), &requestBody); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Bad Request",
		}, nil
	}

	updatedTxn, err := DBService.UpdateTransaction(ctx, requestBody.TransactionId, requestBody.CheckerId, requestBody.Approval)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "",
		}, nil
	}
	responseBody.Txn = *updatedTxn

	bod, err := json.Marshal(responseBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 201,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(bod),
	}, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
