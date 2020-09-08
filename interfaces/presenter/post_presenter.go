package presenter

import (
	"fmt"

	"github.com/slack-go/slack"
)

// postPresenter is presenter for Slack Post
type postPresenter struct {
	client *slack.Client
}

// PostPresenter is a interface
type PostPresenter interface {
	PostMsg(channelID string, msg ...slack.MsgOption) error
	PostBroadCastMsg(msg ...slack.MsgOption) error
}

// NewPostPresenter should be invoked in infrastructure
func NewPostPresenter(client *slack.Client) PostPresenter {
	return &postPresenter{client}
}

//PostMsg posts to a channel
func (pp *postPresenter) PostMsg(channelID string, msgs ...slack.MsgOption) error {
	_, _, err := pp.client.PostMessage(channelID, msgs...)
	if err != nil {
		return fmt.Errorf("Could not parse slash command JSON: %v", err)
	}
	return nil
}

//PostBroadCastMsg post to all channel
func (pp *postPresenter) PostBroadCastMsg(msgs ...slack.MsgOption) error {
	params := slack.GetConversationsParameters{}
	channels, _, err := pp.client.GetConversations(&params)
	if err != nil {
		return fmt.Errorf("GetConversationsParameters %s", err)
	}
	for _, channel := range channels {
		fmt.Printf("Post message to : %v\n", channel.Name)
		_, _, err = pp.client.PostMessage(channel.ID, msgs...)
		if err != nil && err.Error() != "not_in_channel" { // Ignore only this err
			return fmt.Errorf("Could not post message: %v", err)
		}
	}
	return nil
}
