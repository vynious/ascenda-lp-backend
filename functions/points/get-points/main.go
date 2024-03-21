package main

import (
	"context"
	"encoding/json"
	"github.com/vynious/ascenda-lp-backend/types"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/vynious/ascenda-lp-backend/db"
	"gorm.io/gorm"
)

var (
	DBService *db.DBService
	RDSClient *rds.Client
	err       error
)

func init() {
	log.Printf("INIT")
	DBService, err = db.SpawnDBService()
}

func main() {
	// we are simulating a lambda behind an ApiGatewayV2
	lambda.Start(handler)

	defer DBService.CloseConn()
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Printf(request.Body)

	req := types.GetPointsByUserRequestBody{}
	json.Unmarshal([]byte(request.Body), &req)

	conn := DBService.Conn
	var pointsRecords []types.Points

	var res *gorm.DB
	if req.UserId != "" {
		res = conn.Where("user_id = ?", req.UserId).First(&pointsRecords)
	} else {
		res = conn.Find(&pointsRecords)
	}

	if res.Error != nil {
		log.Printf("Database error %s", res.Error)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	}

	if res.RowsAffected == 0 {
		// Return 404 response if no points records are found
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 404,
			Body:       "No points records found",
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
