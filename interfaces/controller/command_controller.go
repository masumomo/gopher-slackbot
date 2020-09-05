package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/masumomo/gopher-slackbot/usecase"
	"github.com/slack-go/slack"
)

// CommandController is controller for Slack Command
type CommandController struct {
	commandInteractor *usecase.CommandInteractor
	api               *slack.Client
	token             string
	verifytoken       string
}

// NewCommandController should be invoked in infrastructure
func NewCommandController(ci *usecase.CommandInteractor) *CommandController {

	return &CommandController{
		commandInteractor: ci,
		api:               slack.New(token),
		token:             os.Getenv("SLACK_BOT_TOKEN"),
		verifytoken:       os.Getenv("SLACK_VERIFY_TOKEN"),
	}
}

//CommandHandler is endpoint for `/commands`
func (cc *CommandController) CommandHandler(w http.ResponseWriter, r *http.Request) {

	sl, err := slack.SlashCommandParse(r)
	if err != nil {
		fmt.Printf("Could not parse slash command JSON: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch sl.Command {
	case "/echo":
		msg := slack.MsgOptionText(sl.Text, false)
		_, _, err = cc.api.PostMessage(sl.ChannelID, msg)
		if err != nil {
			fmt.Printf("Could not post message: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "/echo_broadcast":
		msg := slack.MsgOptionText(sl.Text, false)
		params := slack.GetConversationsParameters{}
		channels, _, err := cc.api.GetConversations(&params)
		if err != nil {
			fmt.Printf("GetConversationsParameters %s\n", err)
			return
		}
		for _, channel := range channels {
			fmt.Printf("Post message to : %v\n", channel.Name)
			_, _, err = cc.api.PostMessage(channel.ID, msg)
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
	case "/schedule":
		// date := strings.TrimSpace(sl.Text)
		msg := slack.MsgOptionText(sl.Text, false)
		params := slack.GetConversationsParameters{}
		channels, _, err := cc.api.GetConversations(&params)
		if err != nil {
			fmt.Printf("GetConversationsParameters %s\n", err)
			return
		}
		for _, channel := range channels {
			fmt.Printf("Post message to : %v\n", channel.Name)
			_, _, err = cc.api.PostMessage(channel.ID, msg)
			if err != nil && err.Error() != "not_in_channel" { // Ignore only this err
				fmt.Printf("Could not post message: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

	default:
		fmt.Printf("This command is not supported : %v\n", sl.Command)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
