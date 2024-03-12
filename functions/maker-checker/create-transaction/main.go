package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/vynious/ascenda-lp-backend/db"
	"github.com/vynious/ascenda-lp-backend/emailer"
	makerchecker "github.com/vynious/ascenda-lp-backend/types/maker-checker"
	"log"
)

var (
	DBService    *db.DBService
	requestBody  makerchecker.CreateTransactionBody
	responseBody makerchecker.CreateMakerResponseBody
	err          error
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env")
	}

	// Initialise global variable DBService tied to Aurora
	DBService, err = db.SpawnDBService()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func LambdaHandler(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	defer DBService.CloseConn()

	/*
		create transaction entry in db => connect to db db
		get checkers based on makers => connect to db to get
		send message through emailer
	*/

	if err := json.Unmarshal([]byte(req.Body), &requestBody); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Bad Request",
		}, nil
	}

	switch requestBody.Action.ResourceType {
	case "User":
		txn, err := DBService.CreateTransaction(ctx, &requestBody)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "",
			}, nil
		}
		responseBody.Txn = *txn
	case "Point":
		txn, err := DBService.CreateTransaction(ctx, &requestBody)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "",
			}, nil
		}
		responseBody.Txn = *txn
	}

	if err = emailer.EmailCheckers(ctx, "<makerId>"); err != nil {
		log.Println(err.Error())
	}

	bod, err := json.Marshal(requestBody)
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
