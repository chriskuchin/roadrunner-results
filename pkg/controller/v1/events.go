package v1

import (
	"fmt"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type eventsResources struct{}

func (rs eventsResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Options("/", Cors)
	r.Get("/", rs.listEvents)
	r.Post("/", rs.addEvent)

	r.Route("/{eventID}", func(r chi.Router) {
		r.Get("/", rs.getEvent)
	})
	return r
}

func (rs eventsResources) getEvent(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("%v", r.Context().Value(util.RaceID))
	render.NoContent(w, r)
}

func (rs eventsResources) listEvents(w http.ResponseWriter, r *http.Request) {
	log.Info().Str("raceID", util.GetRaceIDFromContext(r.Context())).Send()
	if util.GetRaceIDFromContext(r.Context()) == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing RaceID"))
		return
	}

	results, err := services.GetEventsServiceInstance().GetRaceEvents(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unknown error please try again"))
		return
	}
	render.JSON(w, r, results)
}

func (rs eventsResources) addEvent(w http.ResponseWriter, r *http.Request) {
	eventRequest := &EventRequest{}
	if err := render.DecodeJSON(r.Body, eventRequest); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": fmt.Sprintf("%s", err),
		})
		return
	}

	render.NoContent(w, r)
}

type EventRequest struct {
	Description string `json:"description"`
	Distance    int    `json:"distance"`
}
