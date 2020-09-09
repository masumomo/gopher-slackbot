package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/masumomo/gopher-slackbot/usecase"
	"github.com/slack-go/slack"
)

// interactionController is controller for Slack Interaction
type interactionController struct {
	interactionUsecase usecase.InteractionUsecase
}

// InteractionController is controller for Slack Interaction
type InteractionController interface {
	HandleInteraction(r *http.Request) error
}

// NewInteractionController should be invoked in infrastructure
func NewInteractionController(ic usecase.InteractionUsecase) InteractionController {
	return &interactionController{ic}
}

//HandleInteraction is endpoint for `/interactions`
func (ic *interactionController) HandleInteraction(r *http.Request) error {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		return fmt.Errorf("Could not parse action response JSON: %v", err)
	}

	if payload.CallbackID != "select_hello_world" {
		return fmt.Errorf("This calback doesn't support : %v", payload.CallbackID)
	}
	log.Println("Call interaction usecase with:", payload)
	ic.interactionUsecase.RcvInteraction(context.Background(), &payload)
	return nil
}
