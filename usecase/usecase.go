package usecase

import (
	"context"

	"github.com/masumomo/gopher-slackbot/domain/model"
)

type EventUseCase interface {
	SaveEvent(ctx context.Context, event *model.Event) error
}

type InteractionUseCase interface {
	SaveInteraction(ctx context.Context, interaction *model.Interaction) error
}

type CommandUseCase interface {
	SaveCommand(ctx context.Context, command *model.Command) error
}

type GoDocUseCase interface {
	SaveGoDoc(ctx context.Context, goDoc *model.GoDoc) error
}
