package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/vynious/ascenda-lp-backend/db"
	makerchecker "github.com/vynious/ascenda-lp-backend/types"
	"github.com/vynious/ascenda-lp-backend/util"
)

var (
	DBService    *db.DBService
	requestBody  makerchecker.CreateTransactionBody
	responseBody makerchecker.CreateMakerResponseBody
	action       makerchecker.MakerAction
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

func CreateTransactionHandler(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	defer DBService.CloseConn()

	role := "product-manager"                   // do I need?
	makerId := req.RequestContext.Identity.User // ?

	/*
		create transaction entry in db => connect to db
		get checkers based on makers => connect to db to get
		send message through emailer
	*/

	if err := json.Unmarshal([]byte(req.Body), &requestBody); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Bad Request",
		}, nil
	}

	action = makerchecker.MakerAction{
		ResourceType: requestBody.ResourceType,
		ActionType:   requestBody.ActionType,
		RequestBody:  requestBody.RequestBody,
		UserId:       requestBody.UserId,
	}

	// Calls DB Service to create transaction
	txn, err := DBService.CreateTransaction(ctx, action, makerId, requestBody.Description)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "",
		}, nil
	}

	responseBody.Txn = *txn

	// Send emails seek checker's approval
	checkersEmail, err := DBService.GetCheckers(ctx, role)
	if err != nil {
		log.Println(err.Error())
	}

	if err = util.EmailCheckers(ctx, requestBody.ResourceType, checkersEmail); err != nil {
		log.Println(err.Error())
	}

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
	lambda.Start(CreateTransactionHandler)
}
