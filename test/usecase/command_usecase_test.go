package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/masumomo/gopher-slackbot/usecase"

	"github.com/golang/mock/gomock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/masumomo/gopher-slackbot/domain/model"
	"github.com/masumomo/gopher-slackbot/domain/repository"
	mock_datastore "github.com/masumomo/gopher-slackbot/mock/infrastructure/datastore"
	mock_presenter "github.com/masumomo/gopher-slackbot/mock/infrastructure/presenter"
	"github.com/masumomo/gopher-slackbot/test/helper_test"
)

func setupCommand(t *testing.T) (usecase.CommandUsecase, sqlmock.Sqlmock) {

	db, mock := mock_datastore.ConnectMockDB()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS command(.*)").WillReturnResult(sqlmock.NewResult(0, 0))
	r := repository.NewCommandRepository(db)

	gomock := gomock.NewController(t)
	p := mock_presenter.NewMockPostPresenter(gomock)
	uc := usecase.NewCommandUsecase(r, p)
	return uc, mock
}
func TestSaveCommand(t *testing.T) {

	cmd := &model.Command{
		CommandName: "test command name",
		Text:        "test text",
		CreatedBy:   "test user id",
		CreatedAt:   time.Now(),
	}
	uc, mock := setupCommand(t)

	mock.ExpectExec("INSERT INTO commands").WithArgs(cmd.CommandName, cmd.Text, cmd.CreatedBy, helper_test.Anytime{}).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := uc.SaveCommand(context.Background(), cmd.CommandName, cmd.Text, cmd.CreatedBy); err != nil {
		t.Errorf("error was not expected while inserting event: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
