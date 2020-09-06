package usecase

import (
	"context"

	"github.com/masumomo/gopher-slackbot/domain/model"
	"github.com/masumomo/gopher-slackbot/domain/repository"
)

type EventInteractor struct {
	eventRepo *repository.EventRepository
}

type InteractionInteractor struct {
	interactionRepo *repository.InteractionRepository
}

type CommandInteractor struct {
	commandRepo *repository.CommandRepository
}

func NewEventInteractor(eventRepo *repository.EventRepository) *EventInteractor {
	return &EventInteractor{
		eventRepo: eventRepo,
	}
}

func NewInteractionInteractor(interactionRepo *repository.InteractionRepository) *InteractionInteractor {
	return &InteractionInteractor{
		interactionRepo: interactionRepo,
	}
}

func NewCommandInteractor(commandRepo *repository.CommandRepository) *CommandInteractor {
	return &CommandInteractor{
		commandRepo: commandRepo,
	}
}

func (ei *EventInteractor) SaveEvent(ctx context.Context, eventType string, eventText string, createdBy string) error {
	event := model.NewEvent(eventType, eventText, createdBy)
	err := ei.eventRepo.Save(event)
	if err != nil {
		return err
	}
	return nil
}

func (ci *CommandInteractor) SaveCommand(ctx context.Context, commandName string, commandText string, createdBy string) error {
	command := model.NewCommand(commandName, commandText, createdBy)
	err := ci.commandRepo.Save(command)
	if err != nil {
		return err
	}
	return nil
}

func (ei *InteractionInteractor) SaveInteraction(ctx context.Context, interactionType string, interactionData string, createdBy string) error {
	interaction := model.NewInteraction(interactionType, interactionData, createdBy)
	err := ei.interactionRepo.Save(interaction)
	if err != nil {
		return err
	}
	return nil
}

func (ei *EventInteractor) SaveGodDoc(ctx context.Context, goDocName string, goDocURL string, createdBy string) error {
	goDoc := model.NewGoDoc(goDocName, goDocURL, createdBy)
	err := ei.eventRepo.SaveGoDoc(goDoc)
	if err != nil {
		return err
	}
	return nil
}
