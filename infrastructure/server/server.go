package server

import (
	"log"
	"net/http"

	"github.com/masumomo/gopher-slackbot/domain/repository"
	"github.com/masumomo/gopher-slackbot/interfaces/controller"
	"github.com/masumomo/gopher-slackbot/usecase"
)

//App holds server and Interactor to invoke UseCase
type App struct {
	httpServer             *http.Server
	mux                    *http.ServeMux
	eventInteractor        *usecase.EventInteractor
	interactiontInteractor *usecase.InteractionInteractor
	commandInteractor      *usecase.CommandInteractor
}

//NewApp creates repository and UseCase
func NewApp() *App {
	// db := initDB()

	// eventRepo := dbreposigoty.NewEventRepository(db, viper.GetString("mongo.event_collection")) // to use database
	// interactionRepo := dbreposigoty.NewInteractionRepository(db, viper.GetString("mongo.interaction_collection")) // to use database

	eventRepo := repository.NewEventRepository()
	interactionRepo := repository.NewInteractionRepository()
	commandRepo := repository.NewCommandRepository()

	return &App{
		eventInteractor:        usecase.NewEventInteractor(eventRepo),
		interactiontInteractor: usecase.NewInteractionInteractor(interactionRepo),
		commandInteractor:      usecase.NewCommandInteractor(commandRepo),
		mux:                    http.NewServeMux(),
	}
}

// Run is invoked in main at once
func (app *App) Run(port string) error {
	app.httpServer = &http.Server{
		Addr: ":" + port,
	}
	eventController := controller.NewEventController(app.eventInteractor)
	interactionController := controller.NewInteractionController(app.interactiontInteractor)
	commandController := controller.NewCommandController(app.commandInteractor)

	app.mux.HandleFunc("/events", eventController.EventHandler)
	app.mux.HandleFunc("/interactions", interactionController.InteractionHandler)
	app.mux.HandleFunc("/commands", commandController.CommandHandler)

	if err := app.httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen and serve: %+v", err)
		return err
	}
	return nil
}

// func initDB(){

// }
