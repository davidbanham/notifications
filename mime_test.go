package notifications

import (
	"bytes"
	"log"
	"net/mail"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMime(t *testing.T) {
	data := strings.NewReader("oh hi I am an attachment")
	moreData := strings.NewReader("oh hi I am another totally different attachment")

	email := Email{
		To:      TO_ADDRESS,
		From:    "testrun@takehome.io",
		ReplyTo: "lolwut@takehome.io",
		Text:    "this is the text part of a test run",
		HTML:    "this <i>is the HTML part of a test</i> run",
		Subject: "test run",
		Attachments: []Attachment{
			Attachment{
				ContentType: "text/plain",
				Data:        data,
				Filename:    "test_data.txt",
			},
			Attachment{
				ContentType: "text/plain",
				Data:        moreData,
				Filename:    "more_test_data.txt",
			},
		},
	}

	result, err := email.toMIME()
	assert.Nil(t, err)

	log.Printf("DEBUG result: %+v \n", string(result))

	msg, err := mail.ReadMessage(bytes.NewBuffer(result))
	assert.Nil(t, err)
	log.Printf("DEBUG msg: %+v \n", msg)
}
