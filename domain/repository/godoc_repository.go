package repository

import (
	"database/sql"
	"fmt"

	"github.com/masumomo/gopher-slackbot/domain/model"
)

type GoDocRepository struct {
	db *sql.DB
}

func NewGoDocRepository(db *sql.DB) *GoDocRepository {
	// If table doesn't exist create table
	_, err := db.Exec(model.GoDoc.CreateTableDDL(model.GoDoc{}))
	if err != nil {
		panic(err)
	}
	return &GoDocRepository{db: db}
}

func (gr GoDocRepository) Save(goDoc *model.GoDoc) error {
	result, err := gr.db.Exec(fmt.Sprintf("INSERT INTO goDocs(name, url, asked_by,created_by, created_at) VALUES ('%s','%s','%s','%s')", goDoc.Name, goDoc.URL, goDoc.CreatedBy, goDoc.CreatedAt))
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
