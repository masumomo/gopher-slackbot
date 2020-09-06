package repository

import (
	"database/sql"
	"fmt"

	"github.com/masumomo/gopher-slackbot/domain/model"
)

type CommandRepository struct {
	db *sql.DB
}

func NewCommandRepository(db *sql.DB) *CommandRepository {
	// If table doesn't exist create table
	_, err := db.Exec(model.Command.CreateTableDDL(model.Command{}))
	if err != nil {
		panic(err)
	}
	return &CommandRepository{db: db}
}

func (cr CommandRepository) Save(command *model.Command) error {
	result, err := cr.db.Exec("INSERT INTO commands(command_name, text, created_by, created_at) VALUES ($1, $2, $3, $4)", command.CommandName, command.Text, command.CreatedBy, command.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
