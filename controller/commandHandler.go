package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/nlopes/slack/slackevents"
	"github.com/slack-go/slack"
)

var (
	api         *slack.Client
	converter   *md.Converter
	token       string
	verifytoken string
)

func init() {
	converter = md.NewConverter("", true, nil)
	token = os.Getenv("SLACK_BOT_TOKEN")
	verifytoken = os.Getenv("SLACK_VERIFY_TOKEN")
	// generalWebHookURL = os.Getenv("GENERAL_WEB_HOOK_URL")
	// golangWebHookURL = os.Getenv("GOLANG_WEB_HOOK_URL")
	// dm1WebHookURL = os.Getenv("DM1_WEB_HOOK_URL")
	// dm2WebHookURL = os.Getenv("DM2_WEB_HOOK_URL")
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

//WebHookTriggeredByMailHandler is endpoint for `/webhook-triggered-by-mail`
func WebHookTriggeredByMailHandler(w http.ResponseWriter, r *http.Request) {

	type MailFromZapier struct {
		Subject      string `json:"subject"`
		FromName     string `json:"from__name"`
		BodyHTML     string `json:"body_html"`
		BodyMarkdown string
	}

	webHookURL := r.Header.Get("web_hook_url")

	if webHookURL == "" {
		webHookURL = os.Getenv("WEB_HOOK_URL")
		fmt.Printf("Could not get webHookURL. Use default webhook URL: %s\n", webHookURL)
	}

	defer r.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	payload := buf.Bytes()

	var mailFromZapier MailFromZapier
	err := json.Unmarshal(payload, &mailFromZapier)

	if err != nil {
		fmt.Printf("Could not parse json : %v\n", err)
	}
	if err != nil {
		fmt.Printf("Could not decode html : %v\n", err)
	}

	mailFromZapier.BodyMarkdown, err = converter.ConvertString(mailFromZapier.BodyHTML)
	if err != nil {
		fmt.Printf("Could not convert html to markdown : %v\n", err)
	}

	attachment := slack.Attachment{
		Color:    "good",
		Title:    mailFromZapier.Subject + " by " + mailFromZapier.FromName,
		Fallback: "You successfully posted by Incoming Webhook URL!",
		Text:     mailFromZapier.BodyMarkdown,
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}
	err = slack.PostWebhook(webHookURL, &msg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Message successfully sent to channel :%v\n", mailFromZapier.Subject)
}

//WebHookTestHandler is endpoint for `/webhook`
func WebHookTestHandler(w http.ResponseWriter, r *http.Request) {

	webHookURL := r.Header.Get("web_hook_url")

	if webHookURL == "" {
		webHookURL = os.Getenv("WEB_HOOK_URL")
		fmt.Printf("Could not get webHookURL. Use default webhook URL: %s\n", webHookURL)
	}

	attachment := slack.Attachment{
		Color:    "good",
		Title:    "Test webhook",
		Fallback: "You successfully posted by Incoming Webhook URL!",
		Text:     "This is a sentence with some `inline *code*` in it. :smile:\nSee <http://www.foo.com|This message *is* a link>",
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}
	err := slack.PostWebhook(webHookURL, &msg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Message successfully sent to channel!\n")
}
