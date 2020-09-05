package controller

import (
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
