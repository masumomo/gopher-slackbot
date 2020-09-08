package main

import (
	"fmt"
	"os"

	"github.com/masumomo/gopher-slackbot/infrastructure/datastore"
	"github.com/masumomo/gopher-slackbot/infrastructure/web"
	_ "github.com/masumomo/gopher-slackbot/infrastructure/web"
	"github.com/masumomo/gopher-slackbot/register"
	"github.com/slack-go/slack"
)

func main() {
	fmt.Println("[INFO] Server listening")

	db := datastore.ConnectDB()

	port := os.Getenv("PORT")
	verifytoken := os.Getenv("SLACK_VERIFY_TOKEN")
	client := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	app := register.NewApp(db, port, client, verifytoken)

	web.NewRouter(app)

	return
}
