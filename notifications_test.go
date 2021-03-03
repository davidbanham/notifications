package notifications

import (
	"strings"
	"testing"
)

func TestNotificationsLive(t *testing.T) {
	err := SendEmail(Email{
		To:      "david@banham.id.au",
		From:    "testrun@takehome.io",
		ReplyTo: "lolwut@takehome.io",
		Text:    "this is the text part of a test run",
		HTML:    "this <i>is the HTML part of a test</i> run",
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
