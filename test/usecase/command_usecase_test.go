package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/slack-go/slack"

	"github.com/masumomo/gopher-slackbot/usecase"

	"github.com/golang/mock/gomock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/masumomo/gopher-slackbot/domain/model"
	"github.com/masumomo/gopher-slackbot/domain/repository"
	mock_datastore "github.com/masumomo/gopher-slackbot/mock/infrastructure/datastore"
	mock_presenter "github.com/masumomo/gopher-slackbot/mock/infrastructure/presenter"
	test_helper "github.com/masumomo/gopher-slackbot/test/helper"
)

func setupCommand(t *testing.T) (usecase.CommandUsecase, *mock_presenter.MockPostPresenter, sqlmock.Sqlmock, *gomock.Controller) {

	db, mockDB := mock_datastore.ConnectMockDB()
	mockDB.ExpectExec("CREATE TABLE IF NOT EXISTS command(.*)").WillReturnResult(sqlmock.NewResult(0, 0))
	r := repository.NewCommandRepository(db)
	mockCtl := gomock.NewController(t)
	p := mock_presenter.NewMockPostPresenter(mockCtl)
	uc := usecase.NewCommandUsecase(r, p)
	return uc, p, mockDB, mockCtl
}
func TestSaveCommand(t *testing.T) {

	cmd := &model.Command{
		CommandName: "test command name",
		Text:        "test text",
		CreatedBy:   "test user id",
		CreatedAt:   time.Now(),
	}

	uc, _, mockDB, mockCtl := setupCommand(t)
	defer mockCtl.Finish()

	mockDB.ExpectExec("INSERT INTO commands").WithArgs(cmd.CommandName, cmd.Text, cmd.CreatedBy, test_helper.AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := uc.SaveCommand(context.Background(), cmd.CommandName, cmd.Text, cmd.CreatedBy); err != nil {
		t.Errorf("error was not expected while inserting event: %s", err)
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
func TestRcvEchoCommand(t *testing.T) {

	cmd := &model.Command{
		CommandName: "/echo",
		Text:        "test text",
		CreatedBy:   "test user id",
		CreatedAt:   time.Now(),
	}

	slCmd := &slack.SlashCommand{
		Command:   cmd.CommandName,
		ChannelID: "test channel id",
		Text:      cmd.Text,
		UserID:    cmd.CreatedBy,
	}

	uc, mockPost, mockDB, mockCtl := setupCommand(t)
	defer mockCtl.Finish()

	// Set expected result
	mockDB.ExpectExec("INSERT INTO commands").WithArgs(cmd.CommandName, cmd.Text, cmd.CreatedBy, test_helper.AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
	mockPost.EXPECT().PostMsg(slCmd.ChannelID, gomock.AssignableToTypeOf(slack.MsgOptionText(slCmd.Text, false))).Return(nil).Times(1)

	// Call actual method
	if err := uc.RcvCommand(context.Background(), slCmd); err != nil {
		t.Errorf("error was not expected while inserting event: %s", err)
	}
	// Check result
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestRcvEchoBroadcastCommand(t *testing.T) {

	cmd := &model.Command{
		CommandName: "/echo_broadcast",
		Text:        "test text",
		CreatedBy:   "test user id",
		CreatedAt:   time.Now(),
	}

	slCmd := &slack.SlashCommand{
		Command:   cmd.CommandName,
		ChannelID: "test channel id",
		Text:      cmd.Text,
		UserID:    cmd.CreatedBy,
	}

	uc, mockPost, mockDB, mockCtl := setupCommand(t)
	defer mockCtl.Finish()

	// Set expected result
	mockDB.ExpectExec("INSERT INTO commands").WithArgs(cmd.CommandName, cmd.Text, cmd.CreatedBy, test_helper.AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
	mockPost.EXPECT().PostBroadCastMsg(gomock.AssignableToTypeOf(slack.MsgOptionText(slCmd.Text, false))).Return(nil).Times(1)

	// Call actual method
	if err := uc.RcvCommand(context.Background(), slCmd); err != nil {
		t.Errorf("error was not expected while inserting event: %s", err)
	}
	// Check result
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestRcvHelloCommand(t *testing.T) {

	cmd := &model.Command{
		CommandName: "/hello",
		Text:        "Hello, I'm so tired... whatup?",
		CreatedBy:   "test user id",
		CreatedAt:   time.Now(),
	}

	slCmd := &slack.SlashCommand{
		Command:   cmd.CommandName,
		ChannelID: "test channel id",
		Text:      cmd.Text,
		UserID:    cmd.CreatedBy,
	}

	uc, mockPost, mockDB, mockCtl := setupCommand(t)
	defer mockCtl.Finish()

	// Set expected result
	mockDB.ExpectExec("INSERT INTO commands").WithArgs(cmd.CommandName, cmd.Text, cmd.CreatedBy, test_helper.AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
	mockPost.EXPECT().PostMsg(slCmd.ChannelID, gomock.AssignableToTypeOf(slack.MsgOptionText(slCmd.Text, false))).Return(nil).Times(1)

	// Call actual method
	if err := uc.RcvCommand(context.Background(), slCmd); err != nil {
		t.Errorf("error was not expected while inserting event: %s", err)
	}
	// Check result
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
