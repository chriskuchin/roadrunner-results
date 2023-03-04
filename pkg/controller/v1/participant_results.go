package v1

import (
	"net/http"

	"github.com/go-chi/chi"
)

type ParticipantResultsResources struct{}

func (pr *ParticipantResultsResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", handleParticipantResults)

	return r
}

func handleParticipantResults(w http.ResponseWriter, r *http.Request) {

}
