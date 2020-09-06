package repository

import (
	"database/sql"

	"github.com/masumomo/gopher-slackbot/domain/model"
)

type InteractionRepository struct {
	db     *sql.DB
	events map[string]*model.Interaction
}

func NewInteractionRepository(db *sql.DB) *InteractionRepository {
	// If table doesn't exist create table
	// _, err := db.Exec(model.Interaction.CreateTableDDL(model.Interaction{}))
	// if err != nil {
	// 	panic(err)
	// }
	return &InteractionRepository{db: db}
}

func (cr *InteractionRepository) Save(instruction *model.Interaction) error {
	// cr.db.Query(fmt.Sprintf("SELECT * FROM interactions WHERE id = '%s'", id))
	// if err != nil {
	// 	return model.Interaction{}, err
	// }
	// if cr.db.TableName("instructions").
	// 	Where(instruction).First(&instruction).RecordNotFound() {

	// 	repository.db.Create(&instruction)
	// 	return nil
	// }
	return nil
}
