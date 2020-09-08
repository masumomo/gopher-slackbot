package register

import (
	"database/sql"

	"github.com/masumomo/gopher-slackbot/domain/repository"
	"github.com/masumomo/gopher-slackbot/interfaces/controller"
	"github.com/masumomo/gopher-slackbot/interfaces/presenter"
	"github.com/masumomo/gopher-slackbot/usecase"
	"github.com/slack-go/slack"
)

//App holds ServeMux and Usecase to invoke Usecase
type App struct {
	Port                  string
	EventController       controller.EventController
	InteractionController controller.InteractionController
	CommandController     controller.CommandController
}

//NewApp creates controller that contains repositories, presenters and usecases
func NewApp(db *sql.DB, port string, slackClient *slack.Client, verifyToken string) *App {

	// Repository
	eventRepo := repository.NewEventRepository(db)
	interactionRepo := repository.NewInteractionRepository(db)
	commandRepo := repository.NewCommandRepository(db)

	//Presenter
	postPres := presenter.NewPostPresenter(slackClient)

	//Usecase
	eventUsecase := usecase.NewEventUsecase(eventRepo, postPres)
	interactionUsecase := usecase.NewInteractionUsecase(interactionRepo, postPres)
	commandUsecase := usecase.NewCommandUsecase(commandRepo, postPres)

	//Create new app
	return &App{
		Port:                  port,
		EventController:       controller.NewEventController(eventUsecase, verifyToken),
		InteractionController: controller.NewInteractionController(interactionUsecase),
		CommandController:     controller.NewCommandController(commandUsecase),
	}
}
