package repository_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/masumomo/gopher-slackbot/domain/model"
	"github.com/masumomo/gopher-slackbot/domain/repository"
	"github.com/masumomo/gopher-slackbot/mock/infrastructure/mock_datastorey"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (r repository.EventRepository, mock sqlmock.Sqlmock) {
	db, mock := mock_datastorey.ConnectDB()
	r := repository.NewEventRepository(db)
	return r, mock
}
func TestSave(t *testing.T) {

	r, mock := setup(t)

	evt := model.Event{
		EventType: "test event",
		Text:      "message",
		CreatedBy: "testuserID",
		CreatedAt: time.Time()
	}

	mock.ExpectExec("INSERT INTO events").WithArgs(evt.EventType, evt.Text, evt.CreatedBy, evt.CreatedAt).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := r.Save(evt); err != nil {
		t.Errorf("error was not expected while inserting event: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
