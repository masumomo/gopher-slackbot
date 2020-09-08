package register

import (
	"database/sql"

	"github.com/masumomo/gopher-slackbot/domain/repository"
	"github.com/masumomo/gopher-slackbot/interfaces/controller"
	"github.com/masumomo/gopher-slackbot/interfaces/presenter"
	"github.com/slack-go/slack"
)

//App holds ServeMux and UseCase to invoke UseCase
type App struct {
	eventController       *controller.EventController
	interactionController *controller.InteractionController
	commandController     *controller.CommandController
}

//NewApp creates controller that contains repositories, presenters and usecases
func NewApp(db *sql.DB, slackClient *slack.Client, verifyToken string) *App {

	// Repository
	eventRepo := repository.NewEventRepository(db)
	interactionRepo := repository.NewInteractionRepository(db)
	commandRepo := repository.NewCommandRepository(db)

	//Presenter
	postPres := presenter.NewPostPresenter(slackClient)

	//Usecase
	eventUseCase := usecase.NewEventUseCase(eventRepo, postPres)
	interactionUseCase := usecase.NewInteractionUseCase(interactionRepo, postPres)
	commandUseCase := usecase.NewCommandUseCase(commandRepo, postPres)

	//Create new app
	return &App{
		eventController:       controller.NewEventController(eventUseCase, verifyToken),
		interactionController: controller.NewInteractionController(interactionUseCase),
		commandController:     controller.NewCommandController(commandUseCase),
	}
}
