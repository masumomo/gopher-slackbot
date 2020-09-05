package controller

import (
	"os"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/slack-go/slack"
)

var (
	api         *slack.Client
	converter   *md.Converter
	token       string
	verifytoken string
)

var randomMessages = []string{
	"Did you need me?",
	"What's up?",
	"I'm quite tired...",
	"I don't wanna work any more",
}

func init() {
	converter = md.NewConverter("", true, nil)
	token = os.Getenv("SLACK_BOT_TOKEN")
	verifytoken = os.Getenv("SLACK_VERIFY_TOKEN")
	api = slack.New(token)
}
