package web

import (
	"log"
	"net/http"

	"github.com/masumomo/gopher-slackbot/interfaces/controller"
	"github.com/masumomo/gopher-slackbot/register"
)

func newEventRouter(ec controller.EventController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := ec.HandleEvent(r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func newInteractionRouter(ic controller.InteractionController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := ic.HandleInteraction(r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
func newCommandRouter(cc controller.CommandController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := cc.HandleCommand(r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func NewRouter(app *register.App) {
	http.HandleFunc("/commands", newCommandRouter(app.CommandController))
	http.HandleFunc("/events", newCommandRouter(app.CommandController))
	http.HandleFunc("/interactions", newCommandRouter(app.CommandController))
	if err := http.ListenAndServe(":"+app.Port, nil); err != nil {
		log.Fatal(err.Error())
	}
}
