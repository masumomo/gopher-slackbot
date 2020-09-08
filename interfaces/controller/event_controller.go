package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/slack-go/slack/slackevents"
)

// EventController is controller for Slack Event
type eventController struct {
	eventUsecase *usecase.eventUsecase
	verifyToken  string
}

// EventController is interface for Slack Event
type EventController interface {
	HandleEvent(r *http.Request)
}

// NewEventController should be invoked in infrastructure
func NewEventController(eu *usecase.eventUsecase, verifyToken string) EventController {
	return &eventController{eu, verifyToken}
}

//HandleEvent is endpoint for `/events`
func (ec *eventController) HandleEvent(r *http.Request) error {
	defer r.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()

	evt, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: ec.verifyToken}))
	if err != nil {
		return fmt.Errorf("Could not parse event JSON: %v", err)
	}

	if evt.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			return fmt.Errorf("Could not parse event response JSON: %v", err)
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}

	ec.eventUsecase.RcvInteraction(context.Background(), evt)
	return nil
}
