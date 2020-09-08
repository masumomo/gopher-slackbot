package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/masumomo/gopher-slackbot/domain/model"
	"github.com/masumomo/gopher-slackbot/domain/repository"
	"github.com/masumomo/gopher-slackbot/interfaces/presenter"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var randomMessages = []string{
	"Did you need me?",
	"What's up?",
	"I'm quite tired...",
	"I don't wanna work any more",
}

type eventUsecase struct {
	eventRepo *repository.EventRepository
	postPres  presenter.PostPresenter
}

type EventUsecase interface {
	SaveEvent(ctx context.Context, eventType string, eventText string, createdBy string) error
	SaveGoDoc(ctx context.Context, goDocName string, url string, createdBy string) error
	RcvEvent(ctx context.Context, evt *slackevents.EventsAPIEvent) error
}

func NewEventUsecase(eventRepo *repository.EventRepository, eventPres presenter.PostPresenter) EventUsecase {
	return &eventUsecase{eventRepo, eventPres}
}

func (eu *eventUsecase) SaveEvent(ctx context.Context, eventType string, eventText string, createdBy string) error {
	event := model.NewEvent(eventType, eventText, createdBy)
	err := eu.eventRepo.Save(event)
	if err != nil {
		return err
	}
	return nil
}

func (eu *eventUsecase) SaveGoDoc(ctx context.Context, eventType string, eventText string, createdBy string) error {
	goDoc := model.NewGoDoc(eventType, eventText, createdBy)
	err := eu.eventRepo.SaveGoDoc(goDoc)
	if err != nil {
		return err
	}
	return nil
}

func (eu *eventUsecase) RcvEvent(ctx context.Context, evt *slackevents.EventsAPIEvent) error {

	switch evt := evt.InnerEvent.Data.(type) {
	case *slackevents.AppMentionEvent: //normal
		err := eu.SaveEvent(context.Background(), evt.Type, evt.Text, evt.User)
		if err != nil {
			fmt.Printf("Could not save event: %v\n", err)
		}
		return eu.postPres.PostMsg(evt.Channel, slack.MsgOptionText("Yes, hello.", false))
	case *slackevents.MessageEvent: //random
		if evt.BotID != "" { //If it came from bot, ignore
			return nil
		}
		err := eu.SaveEvent(context.Background(), evt.Type, evt.Text, evt.User)
		if err != nil {
			return fmt.Errorf("Could not save event: %v", err)
		}
		includesName, _ := regexp.MatchString("(G|g)opher", evt.Text)
		if err != nil {
			return fmt.Errorf("Regex is bad : %v", err)
		}
		if !includesName {
			fmt.Println(evt.Text, " is not matched")
			return nil
		}
		includesTellMe, _ := regexp.MatchString("(T|t)ell me ", evt.Text)
		if err != nil {
			return fmt.Errorf("Regex is bad : %v", err)
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
					return fmt.Errorf("Regex is bad : %v", err)
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
			msg := "Thank you for asking! Here are documentation of *" + pkg + "." + f + "*\n\n"
			refGolangDoc := "https://golang.org/pkg/" + pkg + "/#" + f
			// refDevDoc := "https://devdocs.io/go/" + pkg + "/index#" + f
			err = eu.SaveGoDoc(context.Background(), pkg+"."+f, refGolangDoc, evt.User)
			if err != nil {
				fmt.Printf("Could not save godoc: %v\n", err)
			}
			return eu.postPres.PostMsg(evt.Channel, slack.MsgOptionText(msg+refGolangDoc+"\n", false))
		} else {
			//Reply normal message
			rand.Seed(time.Now().UnixNano())

			return eu.postPres.PostMsg(evt.Channel, slack.MsgOptionText(randomMessages[rand.Intn(len(randomMessages))], false))
		}
	default:
		return fmt.Errorf("This event is not supported : %v", evt)
	}
}
