package web

import (
	"net/http"

	"github.com/masumomo/gopher-slackbot/interfaces/contoroller"
)

func (ec *contoroller.EventController) EventRouter(w http.ResponseWriter, r *http.Request) {
	if ec.HandleEvent(r); err != nil {
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
