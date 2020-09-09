package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/masumomo/gopher-slackbot/domain/model"
	"github.com/masumomo/gopher-slackbot/domain/repository"
	"github.com/masumomo/gopher-slackbot/interfaces/presenter"
	"github.com/slack-go/slack"
)

// interactionUsecase holds interaction repository and post presenter
type interactionUsecase struct {
	interactionRepo *repository.InteractionRepository
	postPres        presenter.PostPresenter
}

// InteractionUsecase is usecase for slack interaction
type InteractionUsecase interface {
	SaveInteraction(ctx context.Context, interactionType string, action string, createdBy string) error
	RcvInteraction(ctx context.Context, payload *slack.InteractionCallback) error
}

// NewInteractionUsecase returns Interaction usecase usecase
func NewInteractionUsecase(interactionRepo *repository.InteractionRepository, postPres presenter.PostPresenter) InteractionUsecase {
	return &interactionUsecase{interactionRepo, postPres}
}

// SaveInteraction saves Interaction model
func (iu *interactionUsecase) SaveInteraction(ctx context.Context, interactionType string, action string, createdBy string) error {
	interaction := model.NewInteraction(interactionType, action, createdBy)
	log.Println("Save interaction :", interaction)
	err := iu.interactionRepo.Save(interaction)
	if err != nil {
		return err
	}
	return nil
}

// RcvInteraction is for slack interaction
func (iu *interactionUsecase) RcvInteraction(ctx context.Context, payload *slack.InteractionCallback) error {

	err := iu.SaveInteraction(context.Background(), string(payload.Type), payload.CallbackID, payload.User.ID)
	if err != nil {
		log.Printf("Could not save interaction: %v\n", err)
	}

	if payload.Type == slack.InteractionTypeMessageAction {
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
		return iu.postPres.PostMsg(payload.Channel.ID, slack.MsgOptionText("I'll show you Hello world code!", false), slack.MsgOptionAttachments(attachment))
	}

	if payload.Type == slack.InteractionTypeInteractionMessage {

		log.Printf("Message button pressed by user %s with value %s\n", payload.User.Name, payload.ActionCallback.AttachmentActions[0].Value)
		var msg slack.MsgOption
		switch payload.ActionCallback.AttachmentActions[0].Value {
		case "golang":
			log.Println("Hello in Go: ", payload.Value)
			msg = slack.MsgOptionText("I'm grad to hear that! \n```fmt.Print(\"Hello World!\")```", false)
		case "python":
			log.Println("Hello in Python: ", payload.Value)
			msg = slack.MsgOptionText("Whhhhhhy?", false)
		case "javascript":
			log.Println("Hello in JavaScript: ", payload.Value)
			msg = slack.MsgOptionText("I see.. \n```console.log(\"Hello World.\")```", false)
		case "ruby":
			log.Println("Hello in Ruby: ", payload.Value)
			msg = slack.MsgOptionText("Whhhhhhhhhhhhy?", false)
		default:
			log.Println("Something is wrong...: ", payload.Value)
			msg = slack.MsgOptionText("Something is wrong...", false)
		}
		return iu.postPres.PostMsg(payload.Channel.ID, msg)

	}
	return fmt.Errorf("This interaction is not supported : %v", string(payload.Type))
}
