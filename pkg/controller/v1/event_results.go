package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type RecordResultRequestPayload struct {
	Timer string `json:"timer_id"`
	End   int64  `json:"end_ts"`
	Bib   string `json:"bib_number"`
}

func EventResultsRoutes(handler *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.handleEventResults)
	r.Post("/", handler.recordResult)

	return r
}

func (api *Handler) handleEventResults(w http.ResponseWriter, r *http.Request) {
	results, err := services.GetEventResults(r.Context(), api.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, results)
}

func (api *Handler) recordResult(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Send()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var payload RecordResultRequestPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Error().Err(err).Send()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := util.SetDB(r.Context(), api.db)
	if payload.End > 0 {
		services.RecordTimerResult(ctx, payload.End)
	}

	if payload.Bib != "" {
		services.RecordFinisherResult(ctx, api.db, payload.Bib)
	}

	render.Status(r, http.StatusAccepted)
}
