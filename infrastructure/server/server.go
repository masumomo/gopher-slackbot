package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/masumomo/gopher-slackbot/domain/repository"
	"github.com/masumomo/gopher-slackbot/interfaces/controller"
	"github.com/masumomo/gopher-slackbot/usecase"
)

//App holds ServeMux and Interactor to invoke UseCase
type App struct {
	db                     *sql.DB
	mux                    *http.ServeMux
	eventInteractor        *usecase.EventInteractor
	interactiontInteractor *usecase.InteractionInteractor
	commandInteractor      *usecase.CommandInteractor
}

//NewApp creates repository and UseCase
func NewApp(db *sql.DB) *App {

	eventRepo := repository.NewEventRepository(db)
	interactionRepo := repository.NewInteractionRepository(db)
	commandRepo := repository.NewCommandRepository(db)

	return &App{
		db:                     db,
		mux:                    http.NewServeMux(),
		eventInteractor:        usecase.NewEventInteractor(eventRepo),
		interactiontInteractor: usecase.NewInteractionInteractor(interactionRepo),
		commandInteractor:      usecase.NewCommandInteractor(commandRepo),
	}
}

// Run is invoked in main at once
func (app *App) Run(port string) error {
	eventController := controller.NewEventController(app.eventInteractor)
	interactionController := controller.NewInteractionController(app.interactiontInteractor)
	commandController := controller.NewCommandController(app.commandInteractor)

	app.mux.HandleFunc("/events", eventController.EventHandler)
	app.mux.HandleFunc("/interactions", interactionController.InteractionHandler)
	app.mux.HandleFunc("/commands", commandController.CommandHandler)

	if err := http.ListenAndServe(":"+port, app.mux); err != nil {
		log.Fatalf("Failed to listen and serve: %+v", err)
		return err
	}
	return nil
}

// func initDB(){

// }
