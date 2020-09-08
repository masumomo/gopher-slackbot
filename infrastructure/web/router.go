package web

import (
	"log"
	"net/http"

	"github.com/masumomo/gopher-slackbot/interfaces/contoroller"
	"github.com/masumomo/gopher-slackbot/register"
)

func EventRouter(w http.ResponseWriter, r *http.Request) {
	//TODO uhhhh
	if app.ec.HandleEvent(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ic *contoroller.InteractionController) InteractionRouter(w http.ResponseWriter, r *http.Request) {
	if ic.HandleInteraction(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (cc *contoroller.CommandController) CommandRouter(w http.ResponseWriter, r *http.Request) {
	if cc.HandleCommand(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func NewRouter(app *register.App) {
	http.HandleFunc("/events", EventRouter)
	http.HandleFunc("/interactions", app.InteractionController.InteractionRouter)
	http.HandleFunc("/commands", app.CommandController.CommandRouter)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to listen and serve: %+v", err)
		return
	}
}
