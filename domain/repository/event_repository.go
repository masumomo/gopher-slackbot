package repository

import (
	"database/sql"
	"fmt"

	"github.com/masumomo/gopher-slackbot/domain/model"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	// If table doesn't exist create table
	_, err := db.Exec(model.Event.CreateTableDDL(model.Event{}))
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(model.GoDoc.CreateTableDDL(model.GoDoc{}))
	if err != nil {
		panic(err)
	}
	return &EventRepository{db: db}
}

func (er EventRepository) Save(event *model.Event) error {
	result, err := er.db.Exec("INSERT INTO events(event_type, text, created_by, created_at) VALUES ($1, $2, $3, $4)", event.EventType, event.Text, event.CreatedBy, event.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

func (er EventRepository) SaveGoDoc(goDoc *model.GoDoc) error {
	result, err := er.db.Exec("INSERT INTO go_docs(name, url, created_by, created_at) VALUES ($1, $2, $3, $4)", goDoc.Name, goDoc.URL, goDoc.CreatedBy, goDoc.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
