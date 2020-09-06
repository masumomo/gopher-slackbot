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
	app := server.NewApp(db)

	port := os.Getenv("PORT")

	if err := app.Run(port); err != nil {
		log.Fatalf("%s", err.Error())
	}

}
