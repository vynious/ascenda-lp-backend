package util

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	aws2 "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	ses2 "github.com/aws/aws-sdk-go/service/ses"
)

// TODO: only aws-v1 can work??
func EmailCheckers(ctx context.Context, actionType string, checkersEmail []string) error {

	sender := "smucomedy@gmail.com"

	sess, err := session.NewSession(&aws2.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("SES_ACCESS_KEY_ID"), os.Getenv("SES_ACCESS_SECRET_KEY"), ""),
	},)
	if err != nil {
		return fmt.Errorf("failed to start sess: %v", err)
	}


	svc := ses2.New(sess)

	body := `Dear Checker,
	You have a pending transaction for approval.
	Please login to view.
		`

	for _, email := range checkersEmail {
		input := &ses2.SendEmailInput{
			Destination: &ses2.Destination{
				ToAddresses: []*string{
					aws.String(email),
				},
			},
			Message: &ses2.Message{
				Body: &ses2.Body{
					Text: &ses2.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(body),
					},
				},
				Subject: &ses2.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(fmt.Sprintf("[Action Required] %s Request", actionType)),
				},
			},
			Source: aws.String(sender),
		}
		_, err = svc.SendEmail(input)

		if err != nil {
			log.Printf("failed send email to %v due to %v", email, err)
		}
	}
	return nil
}

// VerifyEmail sends a verification email to the target address to add their verify their identity for receiving emails.
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

func VerifyEmail(ctx context.Context, email string) error {
	return nil
}

