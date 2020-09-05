package repository

import (
	"github.com/masumomo/gopher-slackbot/domain/model"
)

type CommandRepository struct {
	events map[string]*model.Event
}

func NewCommandRepository() *CommandRepository {
	return &CommandRepository{}
}

func (er CommandRepository) Save(*model.Command) error {
	return nil
}
