package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/masumomo/gopher-slackbot/domain/model"
	"github.com/masumomo/gopher-slackbot/domain/repository"
	mock_datastore "github.com/masumomo/gopher-slackbot/mock/infrastructure/datastore"
)

func setupCommand(t *testing.T) (*repository.CommandRepository, sqlmock.Sqlmock) {

	db, mock := mock_datastore.ConnectMockDB()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS command(.*)").WillReturnResult(sqlmock.NewResult(0, 0))
	r := repository.NewCommandRepository(db)
	return r, mock
}
func TestSaveCommand(t *testing.T) {

	cmd := &model.Command{
		CommandName: "test command name",
		Text:        "test text",
		CreatedBy:   "test user id",
		CreatedAt:   time.Now(),
	}
	r, mock := setupCommand(t)

	mock.ExpectExec("INSERT INTO commands").WithArgs(cmd.CommandName, cmd.Text, cmd.CreatedBy, cmd.CreatedAt).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := r.Save(cmd); err != nil {
		t.Errorf("error was not expected while inserting event: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
