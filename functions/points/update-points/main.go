package main

import (
	"context"
	"encoding/json"
	"errors"
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
	req := types.UpdatePointsRequestBody{}
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       errors.New("invalid request. malformed request found").Error(),
		}, nil
	}

	if req.ID == nil || req.NewBalance == nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       errors.New("bad request. id or new_balance not found").Error(),
		}, nil
	}
	log.Printf("UpdatePoints %s", *req.ID)

	pointsRecord, err := DB.UpdatePoints(ctx, req)
	if pointsRecord == nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	obj, _ := json.Marshal(pointsRecord)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(obj),
	}, nil
}
