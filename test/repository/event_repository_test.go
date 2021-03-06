package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/masumomo/gopher-slackbot/domain/model"
	"github.com/masumomo/gopher-slackbot/domain/repository"
	mock_datastore "github.com/masumomo/gopher-slackbot/mock/infrastructure/datastore"
)

func setupEvent(t *testing.T) (*repository.EventRepository, sqlmock.Sqlmock) {

	db, mock := mock_datastore.ConnectMockDB()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS events(.*)").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS go_docs(.*)").WillReturnResult(sqlmock.NewResult(0, 0))
	r := repository.NewEventRepository(db)
	return r, mock
}
func TestSaveEvent(t *testing.T) {

	evt := &model.Event{
		EventType: "test event",
		Text:      "message",
		CreatedBy: "testuserID",
		CreatedAt: time.Now(),
	}
	r, mock := setupEvent(t)

	mock.ExpectExec("INSERT INTO events").WithArgs(evt.EventType, evt.Text, evt.CreatedBy, evt.CreatedAt).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := r.Save(evt); err != nil {
		t.Errorf("error was not expected while inserting event: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestSaveGoDoc(t *testing.T) {

	evt := &model.GoDoc{
		Name:      "go packege.func name",
		URL:       "http;//sample.com",
		CreatedBy: "test user id",
		CreatedAt: time.Now(),
	}
	r, mock := setupEvent(t)

	mock.ExpectExec("INSERT INTO go_docs").WithArgs(evt.Name, evt.URL, evt.CreatedBy, evt.CreatedAt).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := r.SaveGoDoc(evt); err != nil {
		t.Errorf("error was not expected while inserting event: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
