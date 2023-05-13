package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func EventsRoutes(handler *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.listEvents)
	r.Post("/", handler.addEvent)

	r.Route("/{eventID}", func(r chi.Router) {
		r.Use(eventCtx)
		r.Get("/", handler.getEvent)
		r.Delete("/", handler.deleteEvent)
		r.Mount("/results", EventResultsRoutes(handler))
		r.Mount("/timers", TimerRoutes(handler))
	})
	return r
}

func (api *Handler) deleteEvent(w http.ResponseWriter, r *http.Request) {
	services.DeleteRaceEvent(r.Context(), api.db)
}

func (api *Handler) getEvent(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("%v", r.Context().Value(util.RaceID))
	render.NoContent(w, r)
}

func (api *Handler) listEvents(w http.ResponseWriter, r *http.Request) {
	if util.GetRaceIDFromContext(r.Context()) == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing RaceID"))
		return
	}

	results, err := services.GetRaceEvents(r.Context(), api.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unknown error please try again"))
		return
	}
	render.JSON(w, r, results)
}

func (api *Handler) addEvent(w http.ResponseWriter, r *http.Request) {
	eventRequest := &EventRequest{}
	if err := render.DecodeJSON(r.Body, eventRequest); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": fmt.Sprintf("%s", err),
		})
		return
	}

	services.AddEvent(r.Context(), api.db, util.GetRaceIDFromContext(r.Context()), eventRequest.Description, eventRequest.Type, eventRequest.Distance)

	render.NoContent(w, r)
}

func eventCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventID := chi.URLParam(r, "eventID")
		if eventID == "" {
			log.Error().Msg("Unable to locate eventID")
		}

		ctx := context.WithValue(r.Context(), util.EventID, eventID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type EventRequest struct {
	Description string `json:"description"`
	Type        string `json:"type"`
	Distance    int    `json:"distance"`
}
