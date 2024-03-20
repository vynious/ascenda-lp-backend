package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/vynious/ascenda-lp-backend/db"
	"github.com/vynious/ascenda-lp-backend/types/points"
)

var (
	DBService *db.DBService
	RDSClient *rds.Client
	err       error
)

func init() {
	log.Printf("INIT")
	DBService, err = db.SpawnDBService()

	defer DBService.CloseConn()
}

func main() {
	// we are simulating a lambda behind an ApiGatewayV2
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Printf(os.Getenv("dbHost"))
	log.Printf(request.Body)

	conn := DBService.Conn
	var pointsRecords []points.Points

	res := conn.Find(&pointsRecords)
	if res.RowsAffected == 0 {
		log.Printf("No points record found %s", res.Error)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 200,
			Body:       "No record found",
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       "Hello",
	}, nil
}
