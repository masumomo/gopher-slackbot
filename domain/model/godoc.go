package model

import "time"

type GoDoc struct {
	Name      string
	URL       string
	CreatedBy string
	CreatedAt time.Time
}

func NewGoDoc(name, url, userName string) *GoDoc {
	return &GoDoc{
		Name:      name,
		URL:       url,
		CreatedBy: userName,
		CreatedAt: time.Now(),
	}
}

func (GoDoc) TableName() string { return "go_docs" }

func (GoDoc) CreateTableDDL() string {
	return `CREATE TABLE IF NOT EXISTS go_docs(` +
		`"id" SERIAL PRIMARY KEY` +
		`,name VARCHAR(255) NOT NULL` +
		`,url VARCHAR(255) NOT NULL` +
		`,created_by VARCHAR(255) NOT NULL` +
		`,created_at TIMESTAMP NOT NULL` +
		`);`
}
