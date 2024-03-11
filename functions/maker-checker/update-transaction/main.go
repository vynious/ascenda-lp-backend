package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	*/

	
    response := events.APIGatewayProxyResponse{}
    return response, nil
}	


func main() {
	lambda.Start(LambdaHandler)
}