package usecase

import (
	"context"

	"github.com/masumomo/gopher-slackbot/domain/models"
)

type EventUseCase interface {
	GetEvents(ctx context.Context, event *models.Event) ([]*models.Event, error)
	DeleteEvent(ctx context.Context, event *models.Event, id string) error
	SaveEvent(ctx context.Context, event *models.Event, id string) error
}

type InteractionUseCase interface {
	GetInteractions(ctx context.Context, interaction *models.Interaction) ([]*models.Interaction, error)
	DeleteInteraction(ctx context.Context, interaction *models.Interaction, id string) error
	SaveInteraction(ctx context.Context, interaction *models.Interaction, id string) error
}

type CommandUseCase interface {
	GetCommands(ctx context.Context, command *models.Command) ([]*models.Command, error)
	DeleteCommand(ctx context.Context, command *models.Command, id string) error
	SaveCommand(ctx context.Context, command *models.Command, id string) error
}
