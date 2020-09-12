package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/masumomo/gopher-slackbot/domain/model"
	"github.com/masumomo/gopher-slackbot/domain/repository"
	mock_datastore "github.com/masumomo/gopher-slackbot/mock/infrastructure/datastore"
)

func setupInteraction(t *testing.T) (*repository.InteractionRepository, sqlmock.Sqlmock) {

	db, mock := mock_datastore.ConnectDB()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS interactions(.*)").WillReturnResult(sqlmock.NewResult(0, 0))
	r := repository.NewInteractionRepository(db)
	return r, mock
}
func TestSaveInteraction(t *testing.T) {

	itr := &model.Interaction{
		InteractionType: "test interaction type",
		Action:          "test action",
		CreatedBy:       "test user id",
		CreatedAt:       time.Now(),
	}
	r, mock := setupInteraction(t)

	mock.ExpectExec("INSERT INTO interactions").WithArgs(itr.InteractionType, itr.Action, itr.CreatedBy, itr.CreatedAt).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := r.Save(itr); err != nil {
		t.Errorf("error was not expected while inserting event: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
