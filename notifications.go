package notifications

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/davidbanham/required_env"
)

var svc *ses.SES
var testMode bool
var debugLogging bool

func init() {
	if os.Getenv("TEST_MOCKS_ON") == "true" {
		testMode = true
		return
	}
	if os.Getenv("NOTIFICATIONS_LOG_LEVEL") == "debug" {
		debugLogging = true
	}
	required_env.Ensure(map[string]string{
		"AWS_ACCESS_KEY_ID":     "",
		"AWS_SECRET_ACCESS_KEY": "",
	})
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Fatal(err)
	}
	svc = ses.New(sess)
}

type Email struct {
	To      string
	From    string
	ReplyTo string
	Text    string
	HTML    string
	Subject string
}

func SendEmail(email Email) (err error) {
	if debugLogging {
		log.Printf("DEBUG notifications email: %+v \n", email)
	}
	if testMode {
		log.Println("INFO notifications TESTMODE dropping email to", email.To, "from", email.From)
		return
	}
	log.Println("INFO notifications sending email to", email.To, "from", email.From)

	body := &ses.Body{
		Text: &ses.Content{
			Data:    aws.String(email.Text),
			Charset: aws.String("UTF8"),
		},
	}

	if email.HTML != "" {
		body.Html = &ses.Content{
			Data:    aws.String(email.HTML),
			Charset: aws.String("UTF8"),
		}
	}

	params := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(email.To),
			},
		},
		Message: &ses.Message{
			Body: body,
			Subject: &ses.Content{
				Data:    aws.String(email.Subject),
				Charset: aws.String("UTF8"),
			},
		},
		Source: aws.String(email.From),
	}

	if email.ReplyTo != "" {
		params.ReplyToAddresses = []*string{
			aws.String(email.ReplyTo),
		}
	}

	_, err = svc.SendEmail(params)

	return
}

func SendRawEmail(data []byte) error {
	input := &ses.SendRawEmailInput{
		FromArn: aws.String(""),
		RawMessage: &ses.RawMessage{
			Data: data,
		},
		ReturnPathArn: aws.String(""),
		Source:        aws.String(""),
		SourceArn:     aws.String(""),
	}

	_, err := svc.SendRawEmail(input)
	return err
}
