package model

type Command struct {
	ID string
}

func NewCommand(id string) *Command {
	return &Command{
		ID: id,
	}
}

func (Command) TableName() string { return "commands" }

// func (Command) CreateTableDDL() string {
// 	return `CREATE TABLE IF NOT EXISTS commands(` +
// 		`"id"   INTEGER  NOT NULL PRIMARY KEY` +
// 		`,name VARCHAR(4) NOT NULL` +
// 		`);`
// }
