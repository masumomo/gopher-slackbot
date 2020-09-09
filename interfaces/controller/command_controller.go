package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/masumomo/gopher-slackbot/usecase"
	"github.com/slack-go/slack"
)

// CommandController is controller for Slack Command
type commandController struct {
	commandUsecase usecase.CommandUsecase
}

// CommandController is controller for Slack Command
type CommandController interface {
	HandleCommand(r *http.Request) error
}

// NewCommandController should be invoked in infrastructure
func NewCommandController(cc usecase.CommandUsecase) CommandController {
	return &commandController{cc}
}

//HandleCommand is endpoint for `/commands`
func (cc *commandController) HandleCommand(r *http.Request) error {
	sl, err := slack.SlashCommandParse(r)
	if err != nil {
		return fmt.Errorf("Could not parse slash JSON: %v", err)
	}
	fmt.Println("Call slash command usecase with:", sl)
	cc.commandUsecase.RcvCommand(context.Background(), &sl)
	return nil
}
