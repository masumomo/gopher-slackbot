package web

import (
	"net/http"

	"github.com/masumomo/gopher-slackbot/interfaces/contoroller"
	"github.com/masumomo/gopher-slackbot/register"
)

func newEventRouter(ec *contoroller.EventController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if ec.HandleEvent(r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func newInteractionRouter(ic *contoroller.InteractionController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if ic.HandleInteraction(r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
func newCommandRouter(cc *contoroller.CommandController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if cc.HandleCommand(r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func NewRouter(app *register.App) {
	http.HandleFunc("/commands", newCommandRouter(app.CommandController))
	http.HandleFunc("/events", newCommandRouter(app.CommandController))
	http.HandleFunc("/interactions", newCommandRouter(app.CommandController))
	 if err := http.ListenAndServe(":"+app.Port, nil) && err != nil{
		log.Fatal(err.Error())
	 }
}
