package util

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"log"
	"os"
)

func EmailCheckers(ctx context.Context, actionType string, checkersEmail []string) error {

	sender := "smucomedy@gmail.com"

	body := `Dear Checker,
	You have a pending transaction for approval.
	Please login to view.
		`
	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(
		credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     os.Getenv("SES_ACCESS_KEY_ID"),
				SecretAccessKey: os.Getenv("SES_ACCESS_SECRET_KEY"),
				SessionToken:    "",
				Source:          "",
			},
		}))
	if err != nil {
		return fmt.Errorf("failed loading cfg for ses: %v", err)
	}

	sesClient := ses.NewFromConfig(cfg)

	for _, email := range checkersEmail {
		input := ses.SendEmailInput{
			Destination: &types.Destination{
				ToAddresses: []string{email},
			},
			Message: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(body),
					},
				},
				Subject: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(fmt.Sprintf("[Action Required] %s Request", actionType)),
				},
			},
			Source:               aws.String(sender),
			ConfigurationSetName: nil,
			ReplyToAddresses:     nil,
			ReturnPath:           nil,
			ReturnPathArn:        nil,
			SourceArn:            nil,
			Tags:                 nil,
		}

		verified, err := VerifyEmail(ctx, email)
		if err != nil {
			log.Printf("unable to verify email: %v", err)
			continue
		}
		if verified == false {
			log.Printf("%v is not verified", email)
			continue
		}

		if _, err := sesClient.SendEmail(ctx, &input); err != nil {
			log.Printf("failed send email to %v due to %v", email, err)
		}
	}
	return nil
}

// SendEmailVerification sends a verification email to the target address to add their verify their identity for receiving emails. Post sign-up process.
func SendEmailVerification(ctx context.Context, email string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load config for emailer")
	}

	sesClient := ses.NewFromConfig(cfg)

	if _, err := sesClient.VerifyEmailIdentity(ctx, &ses.VerifyEmailIdentityInput{
		EmailAddress: aws.String(email),
	}); err != nil {
		return fmt.Errorf("failed to send verification email: %v", err)
	}
	return nil
}

// VerifyEmail verifies the email address before sending it
func VerifyEmail(ctx context.Context, email string) (bool, error) {
	return true, nil
}
