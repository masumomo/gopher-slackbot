package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/masumomo/gopher-slackbot/domain/model"
	"github.com/masumomo/gopher-slackbot/domain/repository"
	mock_datastore "github.com/masumomo/gopher-slackbot/mock/infrastructure/datastore"
)

func setup(t *testing.T) (*repository.EventRepository, sqlmock.Sqlmock) {

	db, mock := mock_datastore.ConnectDB()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS events(.*)").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS go_docs(.*)").WillReturnResult(sqlmock.NewResult(0, 0))
	r := repository.NewEventRepository(db)
	return r, mock
}
func TestSave(t *testing.T) {

	evt := &model.Event{
		EventType: "test event",
		Text:      "message",
		CreatedBy: "testuserID",
		CreatedAt: time.Now(),
	}
	r, mock := setup(t)

	mock.ExpectExec("INSERT INTO events").WithArgs(evt.EventType, evt.Text, evt.CreatedBy, evt.CreatedAt).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := r.Save(evt); err != nil {
		t.Errorf("error was not expected while inserting event: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
