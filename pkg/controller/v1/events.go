package v1

import (
	"fmt"
	"net/http"

	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func HandleEventDelete(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := services.DeleteRaceEvent(r.Context(), db)
		if err != nil {
			log.Error().Err(err).Send()
			apiutil.HandleServerError(err, w, r)
			return
		}
		render.NoContent(w, r)
	}
}

func HandleEventGet(db *sqlx.DB) http.HandlerFunc {
	return apiutil.Unimplemented
}

func HandleEventsList(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if util.GetRaceIDFromContext(r.Context()) == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing RaceID"))
			return
		}

		results, err := services.GetRaceEvents(r.Context(), db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unknown error please try again"))
			return
		}
		render.JSON(w, r, results)
	}
}

func HandleEventsCreate(db *sqlx.DB) http.HandlerFunc {
	type EventRequest struct {
		Description string `json:"description"`
		Type        string `json:"type"`
		Distance    int    `json:"distance"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		eventRequest := &EventRequest{}
		if err := render.DecodeJSON(r.Body, eventRequest); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": fmt.Sprintf("%s", err),
			})
			return
		}

		services.AddEvent(r.Context(), db, util.GetRaceIDFromContext(r.Context()), eventRequest.Description, eventRequest.Type, eventRequest.Distance)

		render.NoContent(w, r)
	}
}
