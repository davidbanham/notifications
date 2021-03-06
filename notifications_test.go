package notifications

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotificationsLive(t *testing.T) {
	err := SendEmail(Email{
		To:      "david@banham.id.au",
		From:    "testrun@takehome.io",
		ReplyTo: "totallynotarealaddress@example.com",
		Text:    "this is the text part of a test run",
		HTML:    "this <i>is the HTML part of a test</i> run",
		Subject: "test run",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotificationsLiveHTMLOnly(t *testing.T) {
	err := SendEmail(Email{
		To:      "david@banham.id.au",
		From:    "testrun@takehome.io",
		ReplyTo: "lolwut@takehome.io",
		HTML:    "this <i>is the HTML, and only, part of a test</i> run",
		Subject: "test run",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotificationsLiveTextOnly(t *testing.T) {
	err := SendEmail(Email{
		To:      "david@banham.id.au",
		From:    "testrun@takehome.io",
		ReplyTo: "lolwut@takehome.io",
		Text:    "this is the text, and only, part of a test run",
		Subject: "test run",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotificationsWithAttachmentsLive(t *testing.T) {
	data := strings.NewReader("oh hi I am an attachment")

	err := SendEmail(Email{
		To:      "david@banham.id.au",
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
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestJSONMarshalAndUnmarshal(t *testing.T) {
	data := strings.NewReader("oh hi I am an attachment")

	email := Email{
		To:      "david@banham.id.au",
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
		},
	}

	result, err := json.Marshal(email)
	assert.Nil(t, err)

	newEmail := Email{}

	assert.Nil(t, json.Unmarshal(result, &newEmail))

	contents, err := ioutil.ReadAll(newEmail.Attachments[0].Data)
	assert.Nil(t, err)
	assert.Equal(t, string(contents), "oh hi I am an attachment")
}
