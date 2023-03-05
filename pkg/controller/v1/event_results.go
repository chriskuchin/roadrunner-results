package v1

import (
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type EventResultsResources struct{}

func (rs EventResultsResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", handleEventResults)

	return r
}

func handleEventResults(w http.ResponseWriter, r *http.Request) {
	results, err := services.GetEventResultsInstance().GetEventResults(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, results)
}
