package main

import (
	"fmt"
	"net/http"

	"time"
	Log "github.com/vynious/ascenda-lp-backend/types"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/macie2"
	"github.com/ip-api/ip-api-go"
)


func CreateLogEntry(log Log) error {
	sess := session.Must(session.NewSession())

	svc := dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		TableName: aws.String("logs"),
		Item: map[string]*dynamodb.AttributeValue{
			"LogID": {
				S: aws.String(log.LogID),
			},
			"UserID": {
				S: aws.String(log.UserID),
			},
			"Type": {
				S: aws.String(log.Type),
			},
			"Action": {
				S: aws.String(log.Action),
			},
			"UserLocation": {
				S: aws.String(log.UserLocation),
			},
			"Timestamp": {
				S: aws.String(log.Timestamp.Format(time.RFC3339)),
			},
			"TTL": {
				S: aws.String(log.TTL),
			},
		},
	}

	_, err := svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func GetLocationFromIP(ip string) (string, error) {
	res, err := ipapi.LookupIP(ip)
	if err != nil {
		return "", err
	}
	location := fmt.Sprintf("%s, %s, %s", res.City, res.RegionName, res.Country)
	return location, nil
}

func filterPIIWithMacie(message string) (string, error) {
	// Create a session using AWS credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a Macie client
	svc := macie2.New(sess)

	// Call Macie to detect PII
	result, err := svc.GetFindings(&macie2.GetFindingsInput{
		FindingIds: []*string{},
	})
	if err != nil {
		return "", fmt.Errorf("failed to get Macie findings: %v", err)
	}

	// Check if Macie identified any PII
	if len(result.Findings) > 0 {
		// Replace PII with a placeholder or mask
		filteredMessage := message
		for _, finding := range result.Findings {
			if strings.Contains(filteredMessage, *finding.Description) {
				filteredMessage = strings.Replace(filteredMessage, *finding.Description, "REDACTED", -1)
			}
		}
		return filteredMessage, nil
	}

	// Return the original message if no PII is found
	return message, nil
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract user ID, action type, and action description from the request
	userID := request.QueryStringParameters["UserID"]
	actionType := request.QueryStringParameters["ActionType"]
	actionDescription := request.QueryStringParameters["ActionDescription"]

	ip := request.RequestContext.Identity.SourceIP

	userLocation, err := GetLocationFromIP(ip)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	filteredDescription, err := FilterPIIWithMacie(actionDescription)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// Create a log entry
	log := Log{
		LogID:         "unique_log_id", 
		UserID:        userID,
		Type:          "user_action",
		Action:        actionDescription,
		UserLocation:  userLocation,
		Timestamp:     time.Now(),
		TTL:           "", 

	// Store the log entry in DynamoDB
	err = CreateLogEntry(log)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Log entry created successfully",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
