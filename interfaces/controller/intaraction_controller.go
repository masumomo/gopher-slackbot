package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/masumomo/gopher-slackbot/usecase"
	"github.com/slack-go/slack"
)

// InteractionController is controller for Slack Interaction
type InteractionController struct {
	interactionInteractor *usecase.InteractionInteractor
}

// NewInteractionController should be invoked in infrastructure
func NewInteractionController(ic *usecase.InteractionInteractor) *InteractionController {
	return &InteractionController{
		interactionInteractor: ic,
	}
}

//InteractionHandler is endpoint for `/interactions`
func (ic *InteractionController) InteractionHandler(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("payload", payload.Type)

	b, err := json.MarshalIndent(payload, "", " ")
	if err == nil {
		s := string(b)
		fmt.Println(s)
	}

	if payload.Type == "shortcut" {
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
		fmt.Println("payload.Channel", payload.Channel)
		fmt.Println("payload.Channel.ID", payload.Channel.ID)
		channelID, timestamp, err := api.PostMessage(payload.Channel.ID, slack.MsgOptionText("I'll show you Hello world code!", false), message)
		if err != nil {
			fmt.Printf("Could not post message: %v\n", err)
			return
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

		_, _, err = api.PostMessage(payload.Channel.ID, msg)
		if err != nil {
			fmt.Printf("Could not post message: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

}
