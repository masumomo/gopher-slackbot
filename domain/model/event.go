package model

import (
	"time"
)

type Event struct {
	EventType string
	Text      string
	CreatedBy string
	CreatedAt time.Time
}

func NewEvent(eventType, text, userName string) *Event {
	return &Event{
		EventType: eventType,
		Text:      text,
		CreatedBy: userName,
		CreatedAt: time.Now(),
	}
}

func (Event) TableName() string { return "events" }

func (Event) CreateTableDDL() string {
	return `CREATE TABLE IF NOT EXISTS events(` +
		`"id" SERIAL PRIMARY KEY` +
		`,event_type VARCHAR(25) NOT NULL` +
		`,"text" TEXT NOT NULL` +
		`,created_by VARCHAR(255) NOT NULL` +
		`,created_at TIMESTAMP NOT NULL` +
		`);`
}
