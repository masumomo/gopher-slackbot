package usecase

import (
	"context"
	"fmt"

	"github.com/masumomo/gopher-slackbot/domain/model"
	"github.com/masumomo/gopher-slackbot/domain/repository"
	"github.com/masumomo/gopher-slackbot/interfaces/presenter"
	"github.com/slack-go/slack"
)

type interactionUseCase struct {
	interactionRepo *repository.InteractionRepository
	postPres        presenter.PostPresenter
}

type InteractionUseCase interface {
	SaveInteraction(ctx context.Context, interactionType string, action string, createdBy string) error
	RcvInteraction(ctx context.Context, payload *slack.InteractionCallback) error
}

func NewInteractionUseCase(interactionRepo *repository.InteractionRepository, postPres presenter.PostPresenter) InteractionUseCase {
	return &interactionUseCase{interactionRepo, postPres}
}

func (iu *interactionUseCase) SaveInteraction(ctx context.Context, interactionType string, action string, createdBy string) error {
	interaction := model.NewInteraction(interactionType, action, createdBy)
	err := iu.interactionRepo.Save(interaction)
	if err != nil {
		return err
	}
	return nil
}

func (iu *interactionUseCase) RcvInteraction(ctx context.Context, payload *slack.InteractionCallback) error {

	err := iu.SaveInteraction(context.Background(), string(payload.Type), payload.ActionID, payload.User.ID)
	if err != nil {
		fmt.Printf("Could not save interaction: %v\n", err)
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
			msg = slack.MsgOptionText("Whhhhhhhhhhhhy?", false)
		default:
			fmt.Println("Something is wrong...: ", payload.Value)
			msg = slack.MsgOptionText("Something is wrong...", false)
		}
		return iu.postPres.PostMsg(payload.Channel.ID, msg)

	}
	return fmt.Errorf("This interaction is not supported : %v", string(payload.Type))
}
