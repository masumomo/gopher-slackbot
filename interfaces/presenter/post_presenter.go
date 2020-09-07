package presenter

import (
	"fmt"
	"net/http"

	"github.com/masumomo/gopher-slackbot/usecase"
	"github.com/slack-go/slack"
)

// PostPresenter is controller for Slack Post
type PostPresenter struct {
	commandUseCase *usecase.PostUseCase
	client         *slack.Client
}

// NewPostPresenter should be invoked in infrastructure
func NewPostPresenter(ci *usecase.PostUseCase, client *slack.Client) *PostPresenter {
	return &PostPresenter{ci, client}
}

//PostEcho is slash command for `/echo`
func (pp *PostPresenter) PostMsg(channelID, msg *slack.MsgOption) error {
	_, _, err = pp.client.PostMessage(channelID, msg)
	if err != nil {
		errMsg := fmt.Sprintf("Could not parse slash command JSON: %v\n", err)
		fmt.Println(errMsg)
		return error.New(errMsg)
	}
	return nil
}

//PostBroadCastMsg is post to slack`
func (pp *PostPresenter) PostBroadCastMsg(msg *slack.MsgOption) error {
	params := slack.GetConversationsParameters{}
	channels, _, err := cp.client.GetConversations(&params)
	if err != nil {
		fmt.Printf("GetConversationsParameters %s\n", err)
		return
	}
	for _, channel := range channels {
		fmt.Printf("Post message to : %v\n", channel.Name)
		_, _, err = cp.client.PostMessage(channel.ID, msg)
		if err != nil && err.Error() != "not_in_channel" { // Ignore only this err
			fmt.Printf("Could not post message: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
