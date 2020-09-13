package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"

	"github.com/masumomo/gopher-slackbot/usecase"

	"github.com/golang/mock/gomock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/masumomo/gopher-slackbot/domain/model"
	"github.com/masumomo/gopher-slackbot/domain/repository"
	mock_datastore "github.com/masumomo/gopher-slackbot/mock/infrastructure/datastore"
	mock_presenter "github.com/masumomo/gopher-slackbot/mock/infrastructure/presenter"
	test_helper "github.com/masumomo/gopher-slackbot/test/helper"
)

func setupEvent(t *testing.T) (usecase.EventUsecase, *mock_presenter.MockPostPresenter, sqlmock.Sqlmock, *gomock.Controller) {

	db, mockDB := mock_datastore.ConnectMockDB()
	mockDB.ExpectExec("CREATE TABLE IF NOT EXISTS events(.*)").WillReturnResult(sqlmock.NewResult(0, 0))
	mockDB.ExpectExec("CREATE TABLE IF NOT EXISTS go_docs(.*)").WillReturnResult(sqlmock.NewResult(0, 0))
	r := repository.NewEventRepository(db)
	mockCtl := gomock.NewController(t)
	p := mock_presenter.NewMockPostPresenter(mockCtl)
	uc := usecase.NewEventUsecase(r, p)
	return uc, p, mockDB, mockCtl
}

func TestSaveEvent(t *testing.T) {

	evt := &model.Event{
		EventType: "test event type",
		Text:      "test text",
		CreatedBy: "test user id",
		CreatedAt: time.Now(),
	}

	uc, _, mockDB, mockCtl := setupEvent(t)
	defer mockCtl.Finish()

	mockDB.ExpectExec("INSERT INTO events").WithArgs(evt.EventType, evt.Text, evt.CreatedBy, test_helper.AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := uc.SaveEvent(context.Background(), evt.EventType, evt.Text, evt.CreatedBy); err != nil {
		t.Errorf("error was not expected while inserting event: %s", err)
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRcvMentionEvent(t *testing.T) {

	evt := &model.Event{
		EventType: "test event type",
		Text:      "test text",
		CreatedBy: "test user id",
		CreatedAt: time.Now(),
	}

	slackMntEvt := &slackevents.AppMentionEvent{
		Type:    evt.EventType,
		Channel: "test channel id",
		Text:    evt.Text,
		User:    evt.CreatedBy,
	}

	slackEvt := &slackevents.EventsAPIEvent{
		InnerEvent: slackevents.EventsAPIInnerEvent{
			Data: slackMntEvt,
		},
	}

	uc, mockPost, mockDB, mockCtl := setupEvent(t)
	defer mockCtl.Finish()

	// Set expected result
	mockDB.ExpectExec("INSERT INTO events").WithArgs(evt.EventType, evt.Text, evt.CreatedBy, test_helper.AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
	mockPost.EXPECT().PostMsg(slackMntEvt.Channel, gomock.AssignableToTypeOf(slack.MsgOptionText(slackMntEvt.Text, false))).Return(nil).Times(1)

	// Call actual method
	if err := uc.RcvEvent(context.Background(), slackEvt); err != nil {
		t.Errorf("error was not expected while inserting event: %s", err)
	}

	// Check result
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRcvNormalReactionEvent(t *testing.T) {

	evt := &model.Event{
		EventType: "test event type",
		Text:      "hello gopcher",
		CreatedBy: "test user id",
		CreatedAt: time.Now(),
	}

	slackMntEvt := &slackevents.MessageEvent{
		Type:    evt.EventType,
		Channel: "test channel id",
		Text:    evt.Text,
		User:    evt.CreatedBy,
	}

	slackEvt := &slackevents.EventsAPIEvent{
		InnerEvent: slackevents.EventsAPIInnerEvent{
			Data: slackMntEvt,
		},
	}

	uc, mockPost, mockDB, mockCtl := setupEvent(t)
	defer mockCtl.Finish()

	// Set expected result
	mockDB.ExpectExec("INSERT INTO events").WithArgs(evt.EventType, evt.Text, evt.CreatedBy, test_helper.AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
	mockPost.EXPECT().PostMsg(slackMntEvt.Channel, gomock.AssignableToTypeOf(slack.MsgOptionText(slackMntEvt.Text, false))).Return(nil).Times(1)

	// Call actual method
	if err := uc.RcvEvent(context.Background(), slackEvt); err != nil {
		t.Errorf("error was not expected while inserting event: %s", err)
	}

	// Check result
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
