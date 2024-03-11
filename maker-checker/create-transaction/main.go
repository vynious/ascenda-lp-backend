package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	makerchecker "github.com/vynious/ascenda-lp-backend/types/maker-checker"
)

var (
// global variables
)

func init() {
	// init global variables
}

func LambdaHandler(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	/*
		check role/user of requested


		req.RequestContext.Identity??

	*/

	var requestBody makerchecker.CreateMakerRequestBody
	if err := json.Unmarshal([]byte(req.Body), &requestBody); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Bad Request",
		}, nil
	}

	response := events.APIGatewayProxyResponse{}
	return response, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
