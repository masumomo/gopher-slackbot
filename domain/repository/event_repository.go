package repository

import model "github.com/masumomo/gopher-slackbot/domain/models"

type EventRepository struct {
	events map[string]*model.Event
}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

func (er EventRepository) Save(*model.Event) error {
	return nil
}