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

// commandUsecase holds event repository and post presenter
type commandUsecase struct {
	commandRepo *repository.CommandRepository
	postPres    presenter.PostPresenter
}

// CommandUsecase is usecase for slack command
type CommandUsecase interface {
	SaveCommand(ctx context.Context, commandName string, commandText string, createdBy string) error
	RcvCommand(ctx context.Context, sl *slack.SlashCommand) error
}

// NewCommandUsecase returns Command usecase usecase
func NewCommandUsecase(commandRepo *repository.CommandRepository, postPres presenter.PostPresenter) CommandUsecase {
	return &commandUsecase{commandRepo, postPres}
}

// SaveCommand saves Command model
func (cu *commandUsecase) SaveCommand(ctx context.Context, commandName string, commandText string, createdBy string) error {
	command := model.NewCommand(commandName, commandText, createdBy)
	log.Println("Save command :", command)
	err := cu.commandRepo.Save(command)
	if err != nil {
		return err
	}
	return nil
}

// RcvCommand is for slack slash command
func (cu *commandUsecase) RcvCommand(ctx context.Context, sl *slack.SlashCommand) error {

	err := cu.SaveCommand(context.Background(), sl.Command, sl.Text, sl.UserID)
	if err != nil {
		log.Printf("Could not save command: %v\n", err)
	}
	switch sl.Command {
	case "/echo":
		return cu.postPres.PostMsg(sl.ChannelID, slack.MsgOptionText(sl.Text, false))
	case "/echo_broadcast":
		return cu.postPres.PostBroadCastMsg(slack.MsgOptionText(sl.Text, false))
	case "/hello":
		return cu.postPres.PostMsg(sl.ChannelID, slack.MsgOptionText("Hello, I'm so tired... whatup?", false))
	default:
		return fmt.Errorf("This command is not supported : %v", sl.Command)
	}
}
