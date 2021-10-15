package email

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/spf13/viper"
	"html/template"
)

const (
	charset           = "UTF-8"
	subject           = "Your Password has been Reset!"
	emailBodyTemplate = `<h1> Hello, {{.FirstName}}!<h1/>
						<p>You've reset your password successfully!</p>
						<p>Your new password is: {{.NewPassword}}</p>
						<p>If you did ask for this change, please contact us!</p>`
)

type Client struct {
	session *session.Session
	sns     *ses.SES
}

func NewEmailClient(sns *ses.SES) Client {
	return Client{sns: sns}
}

func (c Client) SendPasswordResetEmail(user model.User, newPassword string) error {
	sourceEmail := viper.GetString("AWS_SES_SOURCE_EMAIL")
	messageBody, err := getMessageBody(user, newPassword)
	if err != nil {
		return err
	}

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(user.Email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(messageBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sourceEmail),
	}

	_, err = c.sns.SendEmail(input)

	return err
}

func getMessageBody(user model.User, newPassword string) (string, error) {
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
		return "", err
	}

	return messageBody.String(), nil
}
