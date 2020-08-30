package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/masumomo/gopher-slackbot/controller"
)

func main() {
	http.HandleFunc("/events", controller.EventsHandler)
	http.HandleFunc("/interactions", controller.InteractionsHandler)
	http.HandleFunc("/commands", controller.CommandsHandler)
	http.HandleFunc("/golang-weekly-hook", controller.GolangWeeklyHookHandler)
	http.HandleFunc("/webhook", controller.WebHookTestHandler)
	fmt.Println("[INFO] Server listening")
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, nil)
}
