package mock_datastore

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB{
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	return db
}
}
