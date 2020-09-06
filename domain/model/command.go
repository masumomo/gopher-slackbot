package model

import "time"

type Command struct {
	CommandName string
	Text        string
	CreatedBy   string
	CreatedAt   time.Time
}

func NewCommand(commandName, text, userName string) *Command {
	return &Command{
		CommandName: commandName,
		Text:        text,
		CreatedBy:   userName,
		CreatedAt:   time.Now(),
	}
}

func (Command) TableName() string { return "commands" }

func (Command) CreateTableDDL() string {
	return `CREATE TABLE IF NOT EXISTS commands(` +
		`"id" SERIAL PRIMARY KEY` +
		`,command_name VARCHAR(64) NOT NULL` +
		`,"text" TEXT NOT NULL` +
		`,created_by VARCHAR(255) NOT NULL` +
		`,created_at TIMESTAMP NOT NULL` +
		`);`
}
