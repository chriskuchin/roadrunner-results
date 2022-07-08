package v1

import (
	"net/http"

	"github.com/go-chi/chi"
)

// GET /races/{id}
type racesResources struct{}

func (rs racesResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", listRaces)
	r.Post("/", createRace)

	r.Route("/{raceID}", func(r chi.Router) {
		r.Mount("/participants", participantsResources{}.Routes())
		r.Mount("/events", eventsResources{}.Routes())
	})

	return r
}

func listRaces(w http.ResponseWriter, r *http.Request) {

}

func createRace(w http.ResponseWriter, r *http.Request) {

}
