package aws_helpers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func InitCognitoClient() *cognito.CognitoIdentityProvider {
	sess := session.Must(session.NewSession())
	cogClient := cognitoidentityprovider.New(sess, aws.NewConfig().WithRegion("ap-southeast-1"))

	return cogClient
}
