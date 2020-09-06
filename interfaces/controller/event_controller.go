package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/masumomo/gopher-slackbot/usecase"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// EventController is controller for Slack Event
type EventController struct {
	eventInteractor *usecase.EventInteractor
}

// NewEventController should be invoked in infrastructure
func NewEventController(eu *usecase.EventInteractor) *EventController {
	return &EventController{
		eventInteractor: eu,
	}
}

//EventHandler is endpoint for `/events`
func (ec *EventController) EventHandler(w http.ResponseWriter, r *http.Request) {
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
			ec.eventInteractor.SaveEvent(context.Background(), evt.Type, evt.Text, evt.User)
			_, _, err := api.PostMessage(evt.Channel, slack.MsgOptionText("Yes, hello.", false))
			if err != nil {
				fmt.Printf("Could not post message: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		case *slackevents.MessageEvent:
			if evt.BotID != "" { //If it came from bot, ignore
				return
			}
			ec.eventInteractor.SaveEvent(context.Background(), evt.Type, evt.Text, evt.User)
			includesName, _ := regexp.MatchString("(G|g)opher", evt.Text)
			if err != nil {
				fmt.Printf("Regex is bad : %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !includesName {
				fmt.Println(evt.Text, " is not matched")
				return
			}
			includesTellMe, _ := regexp.MatchString("(T|t)ell me ", evt.Text)
			if err != nil {
				fmt.Printf("Regex is bad : %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var (
				pkg string
				f   string
			)
			if includesTellMe {
				words := strings.Split(evt.Text, " ")
				for _, word := range words {
					isSepalatable, _ := regexp.MatchString("^([a-z])+\\.[A-Z]([A-z])+$", word)
					if err != nil {
						fmt.Printf("Regex is bad : %v\n", err)
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					if isSepalatable {
						//Trim and sepalate message
						pkg = strings.Split(word, ".")[0]
						f = strings.Split(word, ".")[1]
						break
					}
				}
			}

			//Create reply
			if pkg != "" && f != "" {
				//Look for doc
				//TODO it should be in usecase
				rand.Seed(time.Now().UnixNano())
				msg := "Thank you for asking! Here are documentation of *" + pkg + "." + f + "*\n\n"
				refGolangDoc := "https://golang.org/pkg/" + pkg + "/#" + f
				// refDevDoc := "https://devdocs.io/go/" + pkg + "/index#" + f
				ec.eventInteractor.SaveGodDoc(context.Background(), evt.Type, pkg+"."+f, refGolangDoc)
				_, _, err := api.PostMessage(evt.Channel, slack.MsgOptionText(msg+refGolangDoc+"\n", false))
				if err != nil {
					fmt.Printf("Could not post message: %v\n", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				//Reply normal message
				rand.Seed(time.Now().UnixNano())
				_, _, err := api.PostMessage(evt.Channel, slack.MsgOptionText(randomMessages[rand.Intn(len(randomMessages))], false))
				if err != nil {
					fmt.Printf("Could not post message: %v\n", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	}
}
