package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type EventResultsResources struct{}

type RecordResultRequestPayload struct {
	End int64 `json:"end_ts"`
}

func (rs EventResultsResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", handleEventResults)
	r.Post("/", recordResult)

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

func recordResult(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Debug().Err(err).Send()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var payload RecordResultRequestPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Debug().Err(err).Send()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Info().Msgf("test: %+v", payload)

	services.GetEventResultsInstance().RecordResult(r.Context(), payload.End)
}
