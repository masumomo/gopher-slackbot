package repository

import (
	"database/sql"

	"github.com/masumomo/gopher-slackbot/domain/model"
)

type CommandRepository struct {
	db       *sql.DB
	commands map[string]*model.Command
}

func NewCommandRepository(db *sql.DB) *CommandRepository {
	return &CommandRepository{db: db}
}

func (cr CommandRepository) Save(*model.Command) error {
	return nil
}
