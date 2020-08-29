package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var api *slack.Client

var token string
var verifytoken string
var golangChannelID string

func main() {

	token = os.Getenv("SLACK_BOT_TOKEN")
	verifytoken = os.Getenv("SLACK_VERIFY_TOKEN")
	golangChannelID = os.Getenv("CHANNEL_ID")

	api = slack.New(token)

	http.HandleFunc("/events", eventsHandler)
	http.HandleFunc("/interactions", interactionsHandler)
	fmt.Println("[INFO] Server listening")
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, nil)
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()

	verifytoken := os.Getenv("SLACK_VERIFY_TOKEN")
	evt, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: verifytoken}))
	if err != nil {
		fmt.Println("err:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if evt.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			fmt.Println("err:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
		fmt.Println("r.Challenge:", r.Challenge)
	}

	if evt.Type == slackevents.CallbackEvent {
		switch evt := evt.InnerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			_, _, err := api.PostMessage(evt.Channel, slack.MsgOptionText("Yes, hello.", false))
			if err != nil {
				fmt.Println("err:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

func interactionsHandler(w http.ResponseWriter, r *http.Request) {

	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v\n", err)
		return
	}

	if payload.CallbackID != "select_hello_world" {
		fmt.Println("This calback doesn't support :", payload.CallbackID)
		http.Error(w, "This calback doesn't support :"+payload.CallbackID, http.StatusInternalServerError)
	}

	if payload.Type == slack.InteractionTypeShortcut {
		attachment := slack.Attachment{
			Pretext:    "Which programming language do you like?",
			Fallback:   "We don't currently support your client",
			CallbackID: "select_hello_world",
			Color:      "#3AA3E3",
			Actions: []slack.AttachmentAction{
				slack.AttachmentAction{
					Name:  "golang",
					Text:  "Golang",
					Type:  "button",
					Value: "golang",
					Style: "primary",
				},
				slack.AttachmentAction{
					Name:  "javascript",
					Text:  "JavaScript",
					Type:  "button",
					Value: "javascript",
				},
				slack.AttachmentAction{
					Name:  "python",
					Text:  "Python",
					Type:  "button",
					Value: "python",
					Style: "danger",
				},
				slack.AttachmentAction{
					Name:  "ruby",
					Text:  "Ruby",
					Type:  "button",
					Value: "ruby",
					Style: "danger",
				},
			},
		}

		message := slack.MsgOptionAttachments(attachment)
		channelID, timestamp, err := api.PostMessage(golangChannelID, slack.MsgOptionText("I'll show you Hello world code!", false), message)
		if err != nil {
			fmt.Printf("Could not send message: %v\n", err)
		}

		fmt.Printf("Message with buttons successfully sent to channel %s at %s\n", channelID, timestamp)
		return
	}

	if payload.Type == slack.InteractionTypeInteractionMessage {
		fmt.Printf("Message button pressed by user %s with value %s\n", payload.User.Name, payload.ActionCallback.AttachmentActions[0].Value)
		var msg slack.MsgOption
		switch payload.ActionCallback.AttachmentActions[0].Value {
		case "golang":
			fmt.Println("Hello in Go: ", payload.Value)
			msg = slack.MsgOptionText("I'm grad to hear that! \n```fmt.Print(\"Hello World!\")```", false)
		case "python":
			fmt.Println("Hello in Python: ", payload.Value)
			msg = slack.MsgOptionText("Whhhhhhy?", false)
		case "javascript":
			fmt.Println("Hello in JavaScript: ", payload.Value)
			msg = slack.MsgOptionText("I see.. \n```console.log(\"Hello World.\")```", false)
		case "ruby":
			fmt.Println("Hello in Ruby: ", payload.Value)
			msg = slack.MsgOptionText("Whhhhhhy?", false)
		default:
			fmt.Println("Something is wrong...: ", payload.Value)
			msg = slack.MsgOptionText("Something is wrong...", false)

		}

		_, _, err = api.PostMessage(golangChannelID, msg)
		if err != nil {
			fmt.Println("err:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

}
