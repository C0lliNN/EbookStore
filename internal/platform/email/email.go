package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/ebookstore/internal/core/auth"
	"github.com/ebookstore/internal/log"
	"github.com/spf13/viper"
)

const (
	charset           = "UTF-8"
	subject           = "Your Password has been Reset!"
	emailBodyTemplate = `<h1> Hello, {{.FirstName}}!<h1/>
						<p>You've reset your password successfully!</p>
						<p>Your new password is: {{.NewPassword}}</p>
						<p>If you did ask for this change, please contact us!</p>`
)

type Email struct {
	Client *ses.Client
}

func NewSESEmailClient(client *ses.Client) *Email {
	return &Email{Client: client}
}

func (e *Email) SendPasswordResetEmail(ctx context.Context, user auth.User, newPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	log.Infof(ctx, "sending password reset email")

	sourceEmail := viper.GetString("AWS_SES_SOURCE_EMAIL")
	messageBody, err := e.getMessageBody(user, newPassword)
	if err != nil {
		return fmt.Errorf("(SendPasswordResetEmail) failed getting email message body: %w", err)
	}

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			CcAddresses: nil,
			ToAddresses: []string{
				user.Email,
			},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String(charset),
					Data:    aws.String(messageBody),
				},
			},
			Subject: &types.Content{
				Charset: aws.String(charset),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sourceEmail),
	}

	if _, err = e.Client.SendEmail(ctx, input); err != nil {
		return fmt.Errorf("(SendPasswordResetEmail) failed sending email")
	}

	return nil
}

func (e *Email) getMessageBody(user auth.User, newPassword string) (string, error) {
	messageBody := bytes.NewBufferString("")

	tmpl := template.Must(template.New("Password Request Template").Parse(emailBodyTemplate))
	err := tmpl.Execute(messageBody, struct {
		FirstName   string
		NewPassword string
	}{
		FirstName:   user.FirstName,
		NewPassword: newPassword,
	})

	if err != nil {
		return "", fmt.Errorf("(getMessageBody) failed parsing template: %w", err)
	}

	return messageBody.String(), nil
}
