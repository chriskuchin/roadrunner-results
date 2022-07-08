package v1

import (
	"net/http"

	"github.com/go-chi/chi"
)

type participantsResources struct{}

// /races/{id}/participants
// POST /races/{id}/participants/csv
// GET /races/{id}/participants
// GET /races/{id}/participants/{id}/results
func (rs participantsResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.listParticipants)
	r.Route("/{participantID}", func(r chi.Router) {
		r.Get("/", rs.getParticipant)
	})

	r.Post("/csv", rs.importParticipantsCSV)

	return r
}

func (rs participantsResources) importParticipantsCSV(w http.ResponseWriter, r *http.Request) {

}

func (rs participantsResources) listParticipants(w http.ResponseWriter, r *http.Request) {

}

func (rs participantsResources) getParticipant(w http.ResponseWriter, r *http.Request) {

}
