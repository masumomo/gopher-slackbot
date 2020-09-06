package usecase

import (
	"context"
)

type UseCase interface {
	SaveEvent(ctx context.Context, eventType string, eventText string, createdBy string) error
	SaveCommand(ctx context.Context, commandName string, commandText string, createdBy string) error
	SaveInteraction(ctx context.Context, interactionType string, action string, createdBy string) error
	SaveGoDoc(ctx context.Context, goDocName string, url string, createdBy string) error
}
