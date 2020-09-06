package repository

import (
	"database/sql"
	"fmt"

	"github.com/masumomo/gopher-slackbot/domain/model"
)

type InteractionRepository struct {
	db *sql.DB
}

func NewInteractionRepository(db *sql.DB) *InteractionRepository {
	// If table doesn't exist create table
	_, err := db.Exec(model.Interaction.CreateTableDDL(model.Interaction{}))
	if err != nil {
		panic(err)
	}
	return &InteractionRepository{db: db}
}

func (ir InteractionRepository) Save(interaction *model.Interaction) error {
	result, err := ir.db.Exec("INSERT INTO interactions(interaction_type, action,created_at, created_at) VALUES ($1, $2, $3, $4)", interaction.InteractionType, interaction.Action, interaction.CreatedBy, interaction.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
