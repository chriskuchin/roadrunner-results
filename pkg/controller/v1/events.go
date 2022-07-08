package v1

import (
	"net/http"

	"github.com/go-chi/chi"
)

type eventsResources struct{}

func (rs eventsResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.listEvents)

	r.Route("/{eventID}", func(r chi.Router) {
		r.Get("/", rs.getEvent)
	})
	return r
}

func (rs eventsResources) getEvent(w http.ResponseWriter, r *http.Request) {

}

func (rs eventsResources) listEvents(w http.ResponseWriter, r *http.Request) {

}
