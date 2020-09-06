package model

import "time"

type Interaction struct {
	InteractionType string
	Action          string
	CreatedBy       string
	CreatedAt       time.Time
}

func NewInteraction(interactionType, action, userName string) *Interaction {
	return &Interaction{
		InteractionType: interactionType,
		Action:          action,
		CreatedBy:       userName,
		CreatedAt:       time.Now(),
	}
}

func (Interaction) TableName() string { return "interactions" }

func (Interaction) CreateTableDDL() string {
	return `CREATE TABLE IF NOT EXISTS interactions(` +
		`"id" SERIAL PRIMARY KEY` +
		`,interaction_type VARCHAR(25) NOT NULL` +
		`,action VARCHAR(64) NOT NULL` +
		`,created_by VARCHAR(255) NOT NULL` +
		`,created_at TIMESTAMP NOT NULL` +
		`);`
}
