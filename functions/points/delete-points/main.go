package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vynious/ascenda-lp-backend/db"
	"github.com/vynious/ascenda-lp-backend/types"
)

var (
	DB  *db.DBService
	err error
)

func init() {
	log.Printf("INIT")
	DB, err = db.SpawnDBService()
}

func main() {
	// we are simulating a lambda behind an ApiGatewayV2
	lambda.Start(handler)

	defer DB.CloseConn()
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	req := types.DeletePointsAccountRequestBody{}
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       errors.New("invalid request. malformed request found").Error(),
		}, nil
	}

	if req.ID == nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       errors.New("bad request. id not found").Error(),
		}, nil
	}
	log.Printf("UpdatePoints %s", *req.ID)

	deleted, err := DB.DeletePointsAccountByID(ctx, *req.ID)
	if !deleted {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	resp := fmt.Sprintf("Points account %s successfully deleted", *req.ID)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       resp,
	}, nil
}
