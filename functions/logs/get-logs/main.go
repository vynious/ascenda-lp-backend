package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/vynious/ascenda-lp-backend/types"
)

var svc *dynamodb.DynamoDB

func init() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials("AKIAUAWJVJ7IQGRU6FGJ", "tsJfzFep6wKQWX+DSOM20xeLQtiZYQdBkhs2IL3N", ""),
	})
	if err != nil {
		print(err)
	}

	svc = dynamodb.New(sess)
}

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	// Get the log table name from Parameter Store
	const logsTable = "logs"

	// Check if ID is provided in query parameters
	id := request.QueryStringParameters["id"]
	if len(id) > 0 {
		return fetchLogByID(id, logsTable)
	}

	// Fetch all logs if no ID provided
	return fetchLogs(request, logsTable)
}

func fetchLogByID(id, tableName string) (events.APIGatewayProxyResponse, error) {
	// Get single log from DynamoDB
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"log_id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := svc.GetItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{}, errors.New("failed to fetch record")
	}

	if result.Item == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Log does not exist",
		}, nil
	}

	item := new(types.Log)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return events.APIGatewayProxyResponse{}, errors.New("failed to unmarshal record")
	}

	// Convert log item to JSON
	body, _ := json.Marshal(item)

	// Return JSON response
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET",
		},
	}, nil
}

func fetchLogs(request events.APIGatewayV2HTTPRequest, tableName string) (events.APIGatewayProxyResponse, error) {
	// Fetch all logs with pagination (limit 100)
	key := request.QueryStringParameters["key"]
	lastEvaluatedKey := make(map[string]*dynamodb.AttributeValue)

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
		Limit:     aws.Int64(int64(100)),
	}

	if len(key) != 0 {
		lastEvaluatedKey["log_id"] = &dynamodb.AttributeValue{
			S: aws.String(key),
		}
		input.ExclusiveStartKey = lastEvaluatedKey
	}

	result, err := svc.Scan(input)
	if err != nil {
		return events.APIGatewayProxyResponse{}, errors.New("failed to fetch record")
	}

	// Unmarshal fetched logs
	var logs []types.Log
	for _, i := range result.Items {
		logItem := new(types.Log)
		err := dynamodbattribute.UnmarshalMap(i, logItem)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		logs = append(logs, *logItem)
	}

	// Construct response body
	responseBody, _ := json.Marshal(logs)

	// Check if there's more data available
	var lastEvaluatedKeyString string
	if len(result.LastEvaluatedKey) != 0 {
		lastEvaluatedKeyString = *result.LastEvaluatedKey["log_id"].S
	}

	// Return JSON response with pagination key if applicable
	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET",
			"Last-Evaluated-Key":           lastEvaluatedKeyString,
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
