package controller

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
)

var (
	converter *md.Converter
)

var randomMessages = []string{
	"Did you need me?",
	"What's up?",
	"I'm quite tired...",
	"I don't wanna work any more",
}

func init() {
	converter = md.NewConverter("", true, nil)
}
