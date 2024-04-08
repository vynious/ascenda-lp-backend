package util

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/vynious/ascenda-lp-backend/types"
)

func CreateLogEntry(bank string, customLog types.Log) error {
	// Specify your AWS credentials and region here
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("DYNAMODB_ACCESS_KEY_ID"), os.Getenv("DYNAMODB_ACCESS_SECRET_KEY"), ""),
	})
	if err != nil {
		return err
	}
	log.Printf("Connected to log db")
	// Create a DynamoDB client
	svc := dynamodb.New(sess)

	// Filter PII in the action field
	filteredAction := filterPII(customLog.Action)

	// Generate a UUID for the log ID
	logID := uuid.New().String()

	// Set the Timestamp field to the current time
	customLog.Timestamp = time.Now().UTC()
	log.Printf("Putting in db")
	input := &dynamodb.PutItemInput{
		TableName: aws.String(fmt.Sprintf("%s_logs", bank)),
		Item: map[string]*dynamodb.AttributeValue{
			"log_id": {
				S: aws.String(logID),
			},
			"UserID": {
				S: aws.String(customLog.UserId),
			},
			"Type": {
				S: aws.String(customLog.Type),
			},
			"Action": {
				S: aws.String(filteredAction), // Redact PII in the action field
			},
			"UserLocation": {
				S: aws.String(customLog.UserLocation),
			},
			"Timestamp": {
				S: aws.String(customLog.Timestamp.Format(time.RFC3339)),
			},
			"TTL": {
				S: aws.String(customLog.TTL),
			},
		},
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("Error creating log entry: %v", err)
		return err
	}
	log.Printf("putted")
	return nil
}

func filterPII(message string) string {
	// Custom logic to redact PII from the message
	// Redact email addresses
	re := regexp.MustCompile(`[\w\.\-]+@[a-zA-Z0-9\-]+\.[a-zA-Z0-9\-\.]+`)
	filteredMessage := re.ReplaceAllString(message, "[REDACTED_EMAIL]")

	// Redact user names
	userNames := []string{"John", "Doe", "Jane", "Smith"} // Add more user names as needed
	for _, name := range userNames {
		filteredMessage = strings.ReplaceAll(filteredMessage, name, "[REDACTED_NAME]")
	}

	return filteredMessage
}
