package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vynious/ascenda-lp-backend/db"
	"github.com/vynious/ascenda-lp-backend/types"
	"github.com/vynious/ascenda-lp-backend/util"
)

var (
	DBService    *db.DBService
	requestBody  types.CreateTransactionBody
	responseBody types.TransactionResponseBody
	action       types.MakerAction
	err          error
)

func init() {
	// Initialise global variable DBService tied to Aurora
	DBService, err = db.SpawnDBService()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func CreateTransactionHandler(ctx context.Context, req *events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	defer DBService.CloseConn()

	role := "product_owner"

	if err := json.Unmarshal([]byte(req.Body), &requestBody); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 404,
			Body:       "Bad Request",
		}, nil
	}

	makerId := requestBody.MakerId

	switch requestBody.Action.ActionType {
	case "UpdatePoints":

		// to check if the request body matches UpdatePointsRequesBody struct
		var updatePointsRequestBody types.UpdatePointsRequestBody
		if err := json.Unmarshal(requestBody.Action.RequestBody, &updatePointsRequestBody); err != nil {
			log.Printf("Error unmarshalling UpdatePointsRequestBody: %v", err)
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 400,
				Headers: map[string]string{
					"Access-Control-Allow-Headers": "Content-Type",
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "POST",
				},
				Body: "Invalid request format for UpdatePoints",
			}, nil
		}

		log.Printf("UpdatePointsRequestBody: %+v", updatePointsRequestBody)

		// convert to json.RawMessage to fit MakerAction struct
		rawJsonBody, _ := json.Marshal(updatePointsRequestBody)

		// recreate the MakerAction struct to store
		updatedMakerCheckerAction := types.MakerAction{
			ActionType:  "UpdatePoints",
			RequestBody: rawJsonBody,
		}

		txn, err := DBService.CreateTransaction(ctx, updatedMakerCheckerAction, makerId)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 500,
				Headers: map[string]string{
					"Access-Control-Allow-Headers": "Content-Type",
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "POST",
				},
				Body: "",
			}, nil
		}
		responseBody.Txn = *txn
	case "UpdateUser":

	default:
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 404,
			Headers: map[string]string{
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "POST",
			},
			Body: "Bad Request",
		}, nil
	}

	// Send emails seek checker's approval (Async)
	checkersEmail, err := DBService.GetCheckers(ctx, role)
	if err != nil {
		log.Println(err.Error())
	}
	if err = util.EmailCheckers(ctx, requestBody.Action.ActionType,
		checkersEmail); err != nil {
		log.Println(err.Error())
	}

	respBod, err := json.Marshal(responseBody)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 201,
			Headers: map[string]string{
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "POST",
			},
			Body: err.Error(),
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 201,
		Headers: map[string]string{
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		},
		Body: string(respBod),
	}, nil
}

func main() {
	lambda.Start(CreateTransactionHandler)
}
