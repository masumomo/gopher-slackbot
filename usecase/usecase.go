package usecase

import (
	"context"

	"github.com/masumomo/gopher-slackbot/domain/model"
)

type EventUseCase interface {
	GetEvents(ctx context.Context, event *model.Event) ([]*model.Event, error)
	DeleteEvent(ctx context.Context, event *model.Event, id string) error
	SaveEvent(ctx context.Context, event *model.Event, id string) error
}

type InteractionUseCase interface {
	GetInteractions(ctx context.Context, interaction *model.Interaction) ([]*model.Interaction, error)
	DeleteInteraction(ctx context.Context, interaction *model.Interaction, id string) error
	SaveInteraction(ctx context.Context, interaction *model.Interaction, id string) error
}

type CommandUseCase interface {
	GetCommands(ctx context.Context, command *model.Command) ([]*model.Command, error)
	DeleteCommand(ctx context.Context, command *model.Command, id string) error
	SaveCommand(ctx context.Context, command *model.Command, id string) error
}
