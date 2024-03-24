package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/vynious/ascenda-lp-backend/types"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vynious/ascenda-lp-backend/db"
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
	var pointsRecords []types.Points

	if request.Body == "" {
		log.Printf("GetPoints")
		pointsRecords, err = DB.GetPoints(ctx)
	} else {
		log.Printf("GetPointsByUser")
		req := types.GetPointsAccountsByUserRequestBody{}
		if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 400,
				Body:       errors.New("invalid request. malformed request found").Error(),
			}, nil
		}

		log.Printf("GetPointsAccountsByUser %s", *req.UserID)
		pointsRecords, err = DB.GetPointsAccountsByUser(ctx, *req.UserID)
	}

	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	if pointsRecords == nil {
		// Return 404 response if no points records are found
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 404,
			Body:       err.Error(),
		}, nil
	}

	obj, err := json.Marshal(pointsRecords)
	if err != nil {
		log.Printf("Failed to parse points records: %v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(obj),
	}, nil
}
