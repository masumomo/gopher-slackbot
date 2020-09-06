package model

type Event struct {
	ID        string
	EventType string
	Data      string
}

func NewEvent(id, eventType, data string) *Event {
	return &Event{
		ID:        id,
		EventType: eventType,
		Data:      data,
	}
}

func (Event) TableName() string { return "events" }

func (Event) CreateTableDDL() string {
	return `CREATE TABLE IF NOT EXISTS events(` +
		`"id" SERIAL PRIMARY KEY` +
		`,event_type VARCHAR(25) NOT NULL` +
		`,data TEXT NOT NULL` +
		`);`
}
