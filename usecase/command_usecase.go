package usecase

import (
	"context"
	"fmt"

	"github.com/masumomo/gopher-slackbot/domain/model"
	"github.com/masumomo/gopher-slackbot/domain/repository"
	"github.com/masumomo/gopher-slackbot/interfaces/presenter"
	"github.com/slack-go/slack"
)

type commandUsecase struct {
	commandRepo *repository.CommandRepository
	postPres    presenter.PostPresenter
}

type CommandUsecase interface {
	SaveCommand(ctx context.Context, commandName string, commandText string, createdBy string) error
	RcvCommand(ctx context.Context, sl *slack.SlashCommand) error
}

func NewCommandUsecase(commandRepo *repository.CommandRepository, postPres presenter.PostPresenter) CommandUsecase {
	return &commandUsecase{commandRepo, postPres}
}

func (cu *commandUsecase) SaveCommand(ctx context.Context, commandName string, commandText string, createdBy string) error {
	command := model.NewCommand(commandName, commandText, createdBy)
	err := cu.commandRepo.Save(command)
	if err != nil {
		return err
	}
	return nil
}

func (cu *commandUsecase) RcvCommand(ctx context.Context, sl *slack.SlashCommand) error {
	err := cu.SaveCommand(context.Background(), sl.Command, sl.Text, sl.UserID)
	if err != nil {
		fmt.Printf("Could not save command: %v\n", err)
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
