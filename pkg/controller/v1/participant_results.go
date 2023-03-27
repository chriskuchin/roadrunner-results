package v1

import (
	"net/http"

	"github.com/go-chi/chi"
)

func ParticipantResultsRoutes(handler *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.handleParticipantResults)

	return r
}

func (api *Handler) handleParticipantResults(w http.ResponseWriter, r *http.Request) {

}
