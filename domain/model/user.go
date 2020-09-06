package model

import "time"

type User struct {
	ID        string
	Email     string
	CreatedAt time.Time
}

func NewUser(id, email string) *User {
	return &User{
		ID:        id,
		Email:     email,
		CreatedAt: time.Now(),
	}
}

func (User) TableName() string { return "users" }

func (User) CreateTableDDL() string {
	return `CREATE TABLE IF NOT EXISTS users(` +
		`"id"   VARCHAR(16) NOT NULL PRIMARY KEY` +
		`,email VARCHAR(255) ` +
		`,created_at TIMESTAMP NOT NULL` +
		`);`
}
