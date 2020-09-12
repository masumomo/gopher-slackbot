package mock_datastore

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

func ConnectMockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	return db, mock
}
