package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/masumomo/gopher-slackbot/infrastructure/datastore"
	// _ "github.com/masumomo/gopher-slackbot/infrastructure/web"
	"github.com/masumomo/gopher-slackbot/register"
	"github.com/slack-go/slack"
)

func main() {
	fmt.Println("[INFO] Server listening")

	db := datastore.ConnectDB()

	port := os.Getenv("PORT")
	verifytoken := os.Getenv("SLACK_VERIFY_TOKEN")
	client := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	app := register.NewApp(db, client, verifytoken)

	http.HandleFunc("/events", app.EventController.EventRouter)
	http.HandleFunc("/interactions", app.InteractionController.interactionRouter)
	http.HandleFunc("/commands", app.CommandController.commandRouter)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to listen and serve: %+v", err)
		return
	}
	return

}
