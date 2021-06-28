package notifications

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	bandname "github.com/davidbanham/bandname_go"
	"github.com/davidbanham/required_env"
	"github.com/stretchr/testify/assert"
)

var TO_ADDRESS string
var runID string

func init() {
	required_env.Ensure(map[string]string{
		"TO_ADDRESS": "",
	})
	TO_ADDRESS = os.Getenv("TO_ADDRESS")
	runID = bandname.Bandname()
}

func TestNotificationsLiveSimple(t *testing.T) {
	err := SendEmail(Email{
		To:      TO_ADDRESS,
		From:    "testrun@takehome.io",
		ReplyTo: "totallynotarealaddress@example.com",
		Text:    "this is the text part of a test run",
		HTML:    "this <i>is the HTML part of a test</i> run",
		Subject: "Simple Test Run" + runID,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotificationsLiveFancy(t *testing.T) {
	err := SendEmail(Email{
		To:      TO_ADDRESS,
		From:    "testrun@takehome.io",
		ReplyTo: "totallynotarealaddress@example.com",
		Text:    "this is the text part of a test run",
		HTML:    "this <i>is the <b>HTML</b> part of a test</i> run. And it has LINKS <a href=\"https://google.com\">https://google.com</a>",
		Subject: "Fancy Test Run" + runID,
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
		Subject: "HTML Only Test Run" + runID,
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
		Subject: "Text Only Test Run" + runID,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotificationsLiveNonASCIISubject(t *testing.T) {
	err := SendEmail(Email{
		To:      TO_ADDRESS,
		From:    "testrun@takehome.io",
		ReplyTo: "lolwut@takehome.io",
		Text:    "this is the text part of a test run",
		HTML:    "this is the <b>HTML</b> part of a test run",
		Subject: "Non-ASCII subject - “–“ - test run" + runID,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotificationsLiveLongSubjectWithNonASCII(t *testing.T) {
	longSubject := fmt.Sprintf(runID + "Long with non-ASCII - “–“ - test This is a seriously long subject line I mean it is just silly what a ridiculous length of string to put in a subject who would do a think like this it is a bloody outrage do you not know that the maximum length of a MIME header is 75 characters and there's all sorts of nonsense we need to do in order to support multiline headers in combination with encoded words so that non-ASCII characters are supported I mean have you even read rfc2047 20 times?")
	err := SendEmail(Email{
		To:      TO_ADDRESS,
		From:    "testrun@takehome.io",
		ReplyTo: "lolwut@takehome.io",
		Text:    "this is the text part of a test run",
		HTML:    "this is the <b>HTML</b> part of a test run",
		Subject: longSubject,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotificationsLiveLongSubject(t *testing.T) {
	longSubject := fmt.Sprintf("This is a seriously long subject line I mean it is just silly what a ridiculous length of string to put in a subject who would do a think like this it is a bloody outrage do you not know that the maximum length of a MIME header is 75 characters and there's all sorts of nonsense we need to do in order to support multiline headers in combination with encoded words so that non-ASCII characters are supported I mean have you even read rfc2047 20 times?")
	err := SendEmail(Email{
		To:      TO_ADDRESS,
		From:    "testrun@takehome.io",
		ReplyTo: "lolwut@takehome.io",
		Text:    "this is the text part of a test run",
		HTML:    "this is the <b>HTML</b> part of a test run",
		Subject: longSubject,
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
		Subject: "Attachments Test Run" + runID,
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
		Subject: "JSON Test Run" + runID,
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
