package emailer

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

func EmailCheckers(ctx context.Context, makerId string) error {

	var checkersEmail []string
	// todo:  get checkers email from makerId

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load config for emailer")
	}
	sesClient := ses.NewFromConfig(cfg)

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: checkersEmail,
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(""),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(""),
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
