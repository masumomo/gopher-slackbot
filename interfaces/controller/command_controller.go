package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/masumomo/gopher-slackbot/usecase"
	"github.com/slack-go/slack"
)

// CommandController is controller for Slack Command
type CommandController struct {
	commandInteractor *usecase.CommandInteractor
	client            *slack.Client
}

// NewCommandController should be invoked in infrastructure
func NewCommandController(ci *usecase.CommandInteractor, client *slack.Client) *CommandController {
	return &CommandController{
		commandInteractor: ci,
		client:            client,
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

	err = cc.commandInteractor.SaveCommand(context.Background(), sl.Command, sl.Text, sl.UserID)
	if err != nil {
		fmt.Printf("Could not save command: %v\n", err)
	}
	switch sl.Command {
	case "/echo":
		msg := slack.MsgOptionText(sl.Text, false)
		_, _, err = cc.client.PostMessage(sl.ChannelID, msg)
		if err != nil {
			fmt.Printf("Could not post message: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "/echo_broadcast":
		msg := slack.MsgOptionText(sl.Text, false)
		params := slack.GetConversationsParameters{}
		channels, _, err := cc.client.GetConversations(&params)
		if err != nil {
			fmt.Printf("GetConversationsParameters %s\n", err)
			return
		}
		for _, channel := range channels {
			fmt.Printf("Post message to : %v\n", channel.Name)
			_, _, err = cc.client.PostMessage(channel.ID, msg)
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
		channels, _, err := cc.client.GetConversations(&params)
		if err != nil {
			fmt.Printf("GetConversationsParameters %s\n", err)
			return
		}
		for _, channel := range channels {
			fmt.Printf("Post message to : %v\n", channel.Name)
			_, _, err = cc.client.PostMessage(channel.ID, msg)
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
