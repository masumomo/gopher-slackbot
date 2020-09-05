package repository

import model "github.com/masumomo/gopher-slackbot/domain/models"

type CommandRepository struct {
	events map[string]*model.Event
}

func NewCommandRepository() *CommandRepository {
	return &CommandRepository{}
}

func (er CommandRepository) Save(*model.Command) error {
	return nil
}
