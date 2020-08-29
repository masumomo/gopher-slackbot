package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/nlopes/slack/slackevents"
	"github.com/slack-go/slack"
)

var (
	api             *slack.Client
	token           string
	verifytoken     string
	golangChannelID string
)

func init() {
	token = os.Getenv("SLACK_BOT_TOKEN")
	verifytoken = os.Getenv("SLACK_VERIFY_TOKEN")
	golangChannelID = os.Getenv("CHANNEL_ID")
	api = slack.New(token)
}

//EventsHandler is endpoint fot `/events`
func EventsHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()

	verifytoken := os.Getenv("SLACK_VERIFY_TOKEN")
	evt, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: verifytoken}))
	if err != nil {
		fmt.Printf("Could not parse event JSON: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if evt.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			fmt.Printf("Could not parse event response JSON: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}

	if evt.Type == slackevents.CallbackEvent {
		switch evt := evt.InnerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			_, _, err := api.PostMessage(evt.Channel, slack.MsgOptionText("Yes, hello.", false))
			if err != nil {
				fmt.Printf("Could not post message: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

//InteractionsHandler is endpoint fot `/interactions`
func InteractionsHandler(w http.ResponseWriter, r *http.Request) {

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
			fmt.Printf("Could not post message: %v\n", err)
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
			fmt.Printf("Could not post message: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

}

//CommandsHandler is endpoint fot `/commands`
func CommandsHandler(w http.ResponseWriter, r *http.Request) {

	sl, err := slack.SlashCommandParse(r)
	fmt.Println("sl", sl)
	if err != nil {
		fmt.Printf("Could not parse slash command JSON: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch sl.Command {
	case "/echo":
		msg := slack.MsgOptionText(sl.Text, false)
		_, _, err = api.PostMessage(sl.ChannelID, msg)
		if err != nil {
			fmt.Printf("Could not post message: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "/echo_broadcast":
		msg := slack.MsgOptionText(sl.Text, false)
		groups, err := api.GetGroups(false)
		if err != nil {
			fmt.Printf("GetGroups %s\n", err)
			return
		}
		for _, group := range groups {
			fmt.Printf("Post message to : %v\n", group.Name)
			_, _, err = api.PostMessage(group.ID, msg)
			if err != nil {
				fmt.Printf("Could not post message: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	case "/hello":
		params := &slack.Msg{Text: "I'm so tired, hello..."}
		b, err := json.Marshal(params)
		if err != nil {
			fmt.Printf("Could not parse param JSON: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	default:
		fmt.Printf("This command is not supported : %v\n", sl.Command)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
