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
	fmt.Println("[INFO] Server listening")
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, nil)
}