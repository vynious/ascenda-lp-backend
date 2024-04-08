package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/vynious/ascenda-lp-backend/types"
	"github.com/vynious/ascenda-lp-backend/util"
)

var svc *dynamodb.DynamoDB

func init() {
	log.Printf("INIT")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("DYNAMODB_ACCESS_KEY_ID"), os.Getenv("DYNAMODB_ACCESS_SECRET_KEY"), ""),
	})
	if err != nil {
		log.Printf("failed to connect to db: " + err.Error())
		print(err)
	}

	svc = dynamodb.New(sess)
}

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	bank, err := util.GetCustomAttributeWithCognito("custom:bank", request.Headers["Authorization"])
	if err == nil {
		log.Printf("failed to get custom:bank from cognito")
	}
	logsTable := fmt.Sprintf("%s_logs", bank)

	ttlStr := request.QueryStringParameters["TTL"]
	if ttlStr != "" {
		ttl, err := strconv.Atoi(ttlStr)
		if err != nil {
			return events.APIGatewayProxyResponse{}, errors.New("invalid TTL value")
		}

		return setTTL(ttl, logsTable)
	}

	// Check if ID is provided in query parameters
	// id := request.QueryStringParameters["id"]
	// if len(id) > 0 {
	// 	return fetchLogByID(id, logsTable)
	// }

	// Fetch all logs if no ID provided
	return fetchLogs(request, logsTable)
}

func fetchLogs(request events.APIGatewayV2HTTPRequest, tableName string) (events.APIGatewayProxyResponse, error) {

	// start
	key := request.QueryStringParameters["key"]
	lastEvaluatedKey := make(map[string]*dynamodb.AttributeValue)
	item := new([]types.Log)
	itemWithKey := new(types.ReturnLogData)

	log.Printf("fetching logs")
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
		log.Printf("failed to fetch logs: " + err.Error())
		return events.APIGatewayProxyResponse{}, errors.New("failed to fetch record")
	}

	// result, err := dynaClient.Scan(input)
	// if err != nil {
	// 	return nil, errors.New(types.ErrorFailedToFetchRecord)
	// }

	for _, i := range result.Items {
		logItem := new(types.Log)
		err := dynamodbattribute.UnmarshalMap(i, logItem)
		if err != nil {
			return events.APIGatewayProxyResponse{}, errors.New("failed to fetch record")
		}
		*item = append(*item, *logItem)
	}

	itemWithKey.Data = *item

	if len(result.LastEvaluatedKey) == 0 {
		responseBody, err := json.Marshal(itemWithKey)
		if err != nil {
			log.Printf("Failed to marshal response: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Internal server error",
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(responseBody),
		}, nil
	}

	itemWithKey.Key = *result.LastEvaluatedKey["log_id"].S

	responseBody, err := json.Marshal(itemWithKey)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}

// var logs []types.Log
// for _, i := range result.Items {
// 	logItem := new(types.Log)
// 	err := dynamodbattribute.UnmarshalMap(i, logItem)
// 	if err != nil {
// 		log.Printf("failed to unmarshal logs: " + err.Error())
// 		return events.APIGatewayProxyResponse{}, err
// 	}
// 	logItem.LogId = *i["log_id"].S
// 	logs = append(logs, *logItem)
// }

// responseBody, err := json.Marshal(logs)
// if err != nil {
// 	log.Printf("Failed to parse logs: %v", err)
// 	return events.APIGatewayProxyResponse{
// 		StatusCode: 400,
// 		Headers: map[string]string{
// 			"Access-Control-Allow-Headers": "Content-Type",
// 			"Access-Control-Allow-Origin":  "*",
// 			"Access-Control-Allow-Methods": "GET",
// 		},
// 		Body: "Internal server error, failed to execute",
// 	}, nil
// }

// return events.APIGatewayProxyResponse{
// 	Body:       string(responseBody),
// 	StatusCode: 200,
// 	Headers: map[string]string{
// 		"Access-Control-Allow-Headers": "Content-Type",
// 		"Access-Control-Allow-Origin":  "*",
// 		"Access-Control-Allow-Methods": "GET",
// 	},
// }, nil

// func FetchLogs(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*types.ReturnLogData, error) {
//	// get all logs with pagination of limit 100
// 	key := req.QueryStringParameters["key"]
// 	lastEvaluatedKey := make(map[string]*dynamodb.AttributeValue)

// 	item := new([]types.Log)
// 	itemWithKey := new(types.ReturnLogData)

// 	input := &dynamodb.ScanInput{
// 		TableName: aws.String(tableName),
// 		Limit:     aws.Int64(int64(100)),
// 	}

// 	if len(key) != 0 {
// 		lastEvaluatedKey["log_id"] = &dynamodb.AttributeValue{
// 			S: aws.String(key),
// 		}
// 		input.ExclusiveStartKey = lastEvaluatedKey
// 	}

// 	result, err := dynaClient.Scan(input)
// 	if err != nil {
// 		return nil, errors.New(types.ErrorFailedToFetchRecord)
// 	}

// 	for _, i := range result.Items {
// 		logItem := new(types.Log)
// 		err := dynamodbattribute.UnmarshalMap(i, logItem)
// 		if err != nil {
// 			return nil, err
// 		}
// 		*item = append(*item, *logItem)
// 	}

// 	itemWithKey.Data = *item

// 	if len(result.LastEvaluatedKey) == 0 {
// 		return itemWithKey, nil
// 	}

// 	itemWithKey.Key = *result.LastEvaluatedKey["log_id"].S

// 	return itemWithKey, nil
// }

// }

func setTTL(ttl int, tableName string) (events.APIGatewayProxyResponse, error) {
	// Check if TTL is already enabled
	descOutput, err := svc.DescribeTimeToLive(&dynamodb.DescribeTimeToLiveInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Printf("failed to describe TTL: %s", err.Error())
		return events.APIGatewayProxyResponse{}, errors.New("failed to describe TTL")
	}

	// If TTL is already enabled, return an error
	if descOutput.TimeToLiveDescription.TimeToLiveStatus != nil &&
		*descOutput.TimeToLiveDescription.TimeToLiveStatus == "ENABLED" {
		log.Printf("TTL is already enabled for table %s", tableName)
		return events.APIGatewayProxyResponse{}, errors.New("TTL is already enabled")
	}

	input := &dynamodb.UpdateTimeToLiveInput{
		TableName: aws.String(tableName),
		TimeToLiveSpecification: &dynamodb.TimeToLiveSpecification{
			AttributeName: aws.String("ttl"),
			Enabled:       aws.Bool(true),
		},
	}

	_, err = svc.UpdateTimeToLive(input)
	if err != nil {
		log.Printf("failed to update TTL: %s", err.Error())
		return events.APIGatewayProxyResponse{}, errors.New("failed to update TTL")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "TTL set successfully",
	}, nil
}

func main() {
	lambda.Start(handler)
}
