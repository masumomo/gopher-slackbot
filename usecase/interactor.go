package usecase

import "github.com/masumomo/gopher-slackbot/domain/repository"

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
