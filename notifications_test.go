package notifications

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/davidbanham/required_env"
	"github.com/stretchr/testify/assert"
)

var TO_ADDRESS string

func init() {
	required_env.Ensure(map[string]string{
		"TO_ADDRESS": "",
	})
	TO_ADDRESS = os.Getenv("TO_ADDRESS")
}

func TestNotificationsLive(t *testing.T) {
	err := SendEmail(Email{
		To:      TO_ADDRESS,
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
		To:      TO_ADDRESS,
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
		To:      TO_ADDRESS,
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
	moreData := strings.NewReader("oh hi I am another totally different attachment")
	bigData := strings.NewReader(`Sometimes attachments are big

and they have multiple lines and stuff

And that should be okay because we can send all kinds of things via email

Even long things with lots of lines`)
	csvData := strings.NewReader(`lol,data,"is fun"
yes,it,is
even,"when there are",spaces`)

	catData, err := os.Open("./cat.png")
	assert.Nil(t, err)

	err = SendEmail(Email{
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
			Attachment{
				ContentType: "text/plain",
				Data:        bigData,
				Filename:    "big_test_data.txt",
			},
			Attachment{
				ContentType: "text/csv",
				Data:        csvData,
				Filename:    "csv_test_data.csv",
			},
			Attachment{
				ContentType: "image/png",
				Data:        catData,
				Filename:    "cat_test_data.png",
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
