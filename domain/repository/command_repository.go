package repository

import "github.com/masumomo/gopher-slackbot/domain/model"

type CommandRepository struct {
	commands map[string]*model.Command
}

func NewCommandRepository() *CommandRepository {
	return &CommandRepository{}
}

func (cr CommandRepository) Save(*model.Command) error {
	return nil
}
