package aws_helpers

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func InitCognitoClient() *cognito.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	cognitoClient := cognito.NewFromConfig(cfg)
	return cognitoClient
}
