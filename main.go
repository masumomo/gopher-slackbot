package main

import (
	"fmt"
	"log"
	"os"

	"github.com/masumomo/gopher-slackbot/infrastructure/datastore"
	"github.com/masumomo/gopher-slackbot/infrastructure/server"
)

func main() {
	fmt.Println("[INFO] Server listening")

	db := datastore.ConnectDB()

	token := os.Getenv("SLACK_BOT_TOKEN")
	verifytoken := os.Getenv("SLACK_VERIFY_TOKEN")

	app := server.NewApp(db)

	port := os.Getenv("PORT")

	if err := app.Run(port, token, verifytoken); err != nil {
		log.Fatalf("%s", err.Error())
	}

}
