package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/nlopes/slack/slackevents"
	"github.com/slack-go/slack"
)

var (
	api         *slack.Client
	converter   *md.Converter
	webHookUrl  string
	token       string
	verifytoken string
)

func init() {
	converter = md.NewConverter("", true, nil)
	token = os.Getenv("SLACK_BOT_TOKEN")
	verifytoken = os.Getenv("SLACK_VERIFY_TOKEN")
	webHookUrl = os.Getenv("WEB_HOOK_URL")
	api = slack.New(token)
}

//EventsHandler is endpoint for `/events`
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

//InteractionsHandler is endpoint for `/interactions`
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
		channelID, timestamp, err := api.PostMessage(payload.Channel.ID, slack.MsgOptionText("I'll show you Hello world code!", false), message)
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

		_, _, err = api.PostMessage(payload.Channel.ID, msg)
		if err != nil {
			fmt.Printf("Could not post message: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

}

//CommandsHandler is endpoint for `/commands`
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
		params := slack.GetConversationsParameters{}
		channels, _, err := api.GetConversations(&params)
		if err != nil {
			fmt.Printf("GetConversationsParameters %s\n", err)
			return
		}
		for _, channel := range channels {
			fmt.Printf("Post message to : %v\n", channel.Name)
			_, _, err = api.PostMessage(channel.ID, msg)
			if err != nil && err.Error() != "not_in_channel" { // Ignore only this err
				fmt.Printf("Could not post message: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
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

//GolangWeeklyHookHandler is endpoint for `/events`
func GolangWeeklyHookHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	htmlBody := buf.String()

	fmt.Println("html ->", htmlBody)
	markdown, err := converter.ConvertString(htmlBody)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("md ->", markdown)

	resp, err := http.Post(webHookUrl, "application/json", strings.NewReader(markdown))
	if err != nil {
		fmt.Printf("Could not post to  slash : %v\n", err)
	}

	fmt.Printf("Message successfully sent to channel! %v\n", resp)
}

//WebHookTestHandler is endpoint for `/webhook`
func WebHookTestHandler(w http.ResponseWriter, r *http.Request) {

	attachment := slack.Attachment{
		Color:         "good",
		Fallback:      "You successfully posted by Incoming Webhook URL!",
		AuthorName:    "slack-go/slack",
		AuthorSubname: "github.com",
		AuthorLink:    "https://github.com/slack-go/slack",
		AuthorIcon:    "https://avatars2.githubusercontent.com/u/652790",
		Text:          "<!channel> All text in Slack uses the same system of escaping: chat messages, direct messages, file comments, etc. :smile:\nSee <https://api.slack.com/docs/message-formatting#linking_to_channels_and_users>",
		Footer:        "slack api",
		FooterIcon:    "https://platform.slack-edge.com/img/default_application_icon.png",
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}
	err := slack.PostWebhook(webHookUrl, &msg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Message successfully sent to channel!\n")
}
