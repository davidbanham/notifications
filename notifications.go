package notifications

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/davidbanham/mailyak/v3"
	"github.com/davidbanham/required_env"
)

var svc *ses.SES
var testMode bool
var debugLogging bool
var tmpl *template.Template
var attachmentTmpl *template.Template

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
	To          string
	From        string
	ReplyTo     string
	Text        string
	HTML        string
	Subject     string
	Attachments []Attachment
}

type Attachment struct {
	ContentType string
	Data        io.Reader
	Filename    string
}

func (attachment *Attachment) MarshalJSON() ([]byte, error) {
	data, err := ioutil.ReadAll(attachment.Data)
	if err != nil {
		return []byte{}, err
	}

	return json.Marshal(&struct {
		ContentType string `json:"content_type"`
		Data        []byte `json:"data"`
		Filename    string `json:"filename"`
	}{
		ContentType: attachment.ContentType,
		Data:        data,
		Filename:    attachment.Filename,
	})
}

func (attachment *Attachment) UnmarshalJSON(data []byte) error {
	inner := struct {
		ContentType string `json:"content_type"`
		Data        []byte `json:"data"`
		Filename    string `json:"filename"`
	}{}

	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}

	attachment.ContentType = inner.ContentType
	attachment.Filename = inner.Filename
	attachment.Data = bytes.NewReader(inner.Data)

	return nil
}

func SendEmail(email Email) error {
	if debugLogging {
		log.Printf("DEBUG notifications email: %+v \n", email)
	}
	if testMode {
		log.Println("INFO notifications TESTMODE dropping email to", email.To, "from", email.From)
		return nil
	}
	log.Println("INFO notifications sending email to", email.To, "from", email.From)

	yak := mailyak.New("", nil)

	yak.To(email.To)
	yak.From(email.From)
	yak.ReplyTo(email.ReplyTo)
	yak.Subject(email.Subject)
	if email.HTML != "" {
		yak.HTML().Set(email.HTML)
	}
	if email.Text != "" {
		yak.Plain().Set(email.Text)
	}

	for _, attachment := range email.Attachments {
		yak.AttachWithMimeType(attachment.Filename, attachment.Data, attachment.ContentType)
	}

	buf, err := yak.MimeBuf()
	if err != nil {
		return err
	}

	return SendRawEmail(buf.Bytes())
}

func SendRawEmail(data []byte) error {
	if testMode {
		log.Println("INFO notifications TESTMODE dropping raw email")
		return nil
	}

	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: data,
		},
	}

	_, err := svc.SendRawEmail(input)
	return err
}
