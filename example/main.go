package main

import (
	"log"
	"strings"

	"github.com/davidbanham/notifications"
)

func main() {
	if err := notifications.SendEmail(notifications.Email{
		To:      "to@example.com",
		From:    "from@example.com",
		ReplyTo: "reply_to@example.com", //optional
		Text:    "this is the text part of a test run",
		HTML:    "this <i>is the HTML part of a test</i> run",
		Subject: "Simple Test Run",
		Attachments: []notifications.Attachment{
			notifications.Attachment{
				ContentType: "text/plain",
				Filename:    "test_data.txt",
				Data:        strings.NewReader("oh hi I am an attachment"),
			},
		},
	}); err != nil {
		log.Fatal(err)
	}
}
