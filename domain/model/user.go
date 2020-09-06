package model

type User struct {
	ID    string
	Email string
}

func NewUser(id, email string) *User {
	return &User{
		ID:    id,
		Email: email,
	}
}

func (User) TableName() string { return "users" }

// func (User) CreateTableDDL() string {
// 	return `CREATE TABLE IF NOT EXISTS users(` +
// 		`"id"   INTEGER  NOT NULL PRIMARY KEY` +
// 		`,email VARCHAR(255) NOT NULL` +
// 		`);`
// }
