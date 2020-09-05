package repository

import "github.com/masumomo/gopher-slackbot/domain/model"

type InteractionRepository struct {
	events map[string]*model.Event
}

func NewInteractionRepository() *InteractionRepository {
	return &InteractionRepository{}
}

func (er InteractionRepository) Save(*model.Interaction) error {
	return nil
}
