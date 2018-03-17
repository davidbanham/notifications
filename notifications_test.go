package notifications

import (
	"testing"
)

func TestNotificationsLive(t *testing.T) {
	err := SendEmail(Email{
		To:      "david@banham.id.au",
		From:    "testrun@takehome.io",
		ReplyTo: "lolwut@takehome.io",
		Text:    "this is a test run",
		HTML:    "this <i>is a test</i> run",
		Subject: "test run",
	})
	if err != nil {
		t.Fatal(err)
	}
}
