package util

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

func EmailCheckers(ctx context.Context, actionType string, checkersEmail []string) error {

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load config for emailer")
	}
	sesClient := ses.NewFromConfig(cfg)

	body := `
		"Dear Checker,
		
		You have a pending transaction for approval.

		Please login to view.
		`

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{"shawn.thiah.2022@scis.smu.edu.sg"},
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
		Source: aws.String(""),
	}
	_, err = sesClient.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("error sending emails to checkers")
	}
	return nil
}
