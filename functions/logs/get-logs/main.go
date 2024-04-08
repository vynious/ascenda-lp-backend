package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

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

var (
	svc     *dynamodb.DynamoDB
	headers = map[string]string{
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET",
	}
)

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
	if err != nil {
		log.Printf("failed to get custom:bank from cognito")
	}
	logsTable := fmt.Sprintf("%s_logs", bank)
	log.Printf("Fetch logs from %s", logsTable)

	return fetchLogs(request, logsTable)
}

func fetchLogs(request events.APIGatewayV2HTTPRequest, tableName string) (events.APIGatewayProxyResponse, error) {
	// Get search criteria from query parameters
	criteria := request.QueryStringParameters["criteria"]
	value := request.QueryStringParameters["value"]
	key := request.QueryStringParameters["key"]

	lastEvaluatedKey := make(map[string]*dynamodb.AttributeValue)
	item := new([]types.Log)
	itemWithKey := new(types.ReturnLogData)

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
		Limit:     aws.Int64(int64(100)),
	}

	// If a search criteria is provided, add it to the filter expression
	if criteria != "" && value != "" {
		expressionAttributeValues := map[string]*dynamodb.AttributeValue{
			":val": {
				S: aws.String(value),
			},
		}

		// Define filter expressions based on search criteria
		filterExpression := ""
		var expressionAttributeNames map[string]*string // Map to hold expression attribute name aliases

		switch criteria {
		case "LogId":
			filterExpression = "contains(log_id, :val)"
		case "UserId":
			filterExpression = "contains(UserID, :val)"
		case "Type":
			filterExpression = "contains(#T, :val)"
			expressionAttributeNames = map[string]*string{"#T": aws.String("Type")} // Define alias for "Type"
		}

		input.FilterExpression = aws.String(filterExpression)
		input.ExpressionAttributeValues = expressionAttributeValues

		// Assign expression attribute name aliases if they exist
		if len(expressionAttributeNames) > 0 {
			input.ExpressionAttributeNames = expressionAttributeNames
		}
	}

	// If a key is provided, set it as the ExclusiveStartKey
	if len(key) != 0 {
		lastEvaluatedKey["log_id"] = &dynamodb.AttributeValue{
			S: aws.String(key),
		}
		input.ExclusiveStartKey = lastEvaluatedKey
	}

	// Execute the scan operation
	result, err := svc.Scan(input)
	if err != nil {
		log.Printf("failed to fetch logs: " + err.Error())
		return events.APIGatewayProxyResponse{}, errors.New("failed to fetch record")
	}

	// Unmarshal the items and add them to the log array
	for _, i := range result.Items {
		logItem := new(types.Log)
		err := dynamodbattribute.UnmarshalMap(i, logItem)
		if err != nil {
			return events.APIGatewayProxyResponse{}, errors.New("failed to fetch record")
		}
		logItem.LogId = *i["log_id"].S // Assign log ID
		*item = append(*item, *logItem)
	}

	// Populate the ReturnLogData structure
	itemWithKey.Data = *item

	// If there are more records to fetch, set the next key
	if len(result.LastEvaluatedKey) != 0 {
		itemWithKey.Key = *result.LastEvaluatedKey["log_id"].S
	}

	// Format and return the response
	response := formatResponse(itemWithKey)
	return response, nil
}

func formatResponse(itemWithKey *types.ReturnLogData) events.APIGatewayProxyResponse {
	responseBody, err := json.Marshal(itemWithKey)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    headers,
			Body:       "Internal server error",
		}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(responseBody),
	}
}

func main() {
	lambda.Start(handler)
}
