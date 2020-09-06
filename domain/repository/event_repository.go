package repository

import (
	"database/sql"
	"fmt"

	"github.com/masumomo/gopher-slackbot/domain/model"
)

type EventRepository struct {
	db     *sql.DB
	events map[string]*model.Event
}

func NewEventRepository(db *sql.DB) *EventRepository {
	// If table doesn't exist create table
	_, err := db.Exec(model.Event.CreateTableDDL(model.Event{}))
	if err != nil {
		panic(err)
	}
	return &EventRepository{db: db}
}

func (er EventRepository) Save(eventType, eventData string) error {
	result, err := er.db.Exec(fmt.Sprintf("INSERT INTO events(event_type, data) VALUES ('%s','%s')", eventType, eventData))
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
