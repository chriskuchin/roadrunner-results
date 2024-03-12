package v1

import (
	"fmt"
	"net/http"

	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

func HandleEventDelete(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := services.DeleteRaceEvent(ctx, db, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx))
		if err != nil {
			apiutil.HandleServerError(err, w, r)
			return
		}
		render.NoContent(w, r)
	}
}

func HandleEventGet(db *sqlx.DB) http.HandlerFunc {
	type response struct {
		EventID     string  `json:"eventId"`
		Description string  `json:"description"`
		Type        string  `json:"type"`
		Distance    float64 `json:"distance"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		event, err := services.GetRaceEvent(ctx, db, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx))
		if err != nil {
			apiutil.HandleServerError(err, w, r)
			return
		}

		render.JSON(w, r, response{
			EventID:     event.EventID,
			Description: event.Description,
			Type:        event.Type,
			Distance:    event.Distance,
		})
	}
}

func HandleEventsList(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if util.GetRaceIDFromContext(ctx) == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing RaceID"))
			return
		}

		results, err := services.GetRaceEvents(ctx, db, util.GetRaceIDFromContext(ctx))
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
		Description string  `json:"description"`
		Type        string  `json:"type"`
		Distance    float64 `json:"distance"`
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
