package v1

import (
	"fmt"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

// GET /races/{id}
type racesResources struct{}

func (rs racesResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Options("/", Cors)
	r.Get("/", listRaces)
	r.Post("/", createRace)

	r.Route("/{raceID}", func(r chi.Router) {
		r.Options("/", Cors)
		r.Delete("/", deleteRace)

		r.Mount("/participants", participantsResources{}.Routes())
		r.Mount("/events", eventsResources{}.Routes())
	})

	return r
}

func Cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}

// GET /races
func listRaces(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	raceService := services.GetRaceServiceInstance()

	races, _ := raceService.ListRaces(ctx)

	render.JSON(w, r, races)
}

// POST /races
func createRace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	raceService := services.GetRaceServiceInstance()

	raceRequest := &RaceRequest{}
	if err := render.DecodeJSON(r.Body, raceRequest); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": fmt.Sprintf("%v", err),
		})
	}

	raceID, err := raceService.CreateRace(ctx, raceRequest.Name)
	if err != nil {
		log.Debug()
		log.Error().Err(err).Send()
		render.Status(r, http.StatusInternalServerError)
		response := map[string]string{
			"error": fmt.Sprintf("%v", err),
		}
		render.JSON(w, r, response)

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{
		"id": raceID,
	})
}

// DELETE .races/{raceID}
func deleteRace(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "raceID")
	ctx := r.Context()
	err := services.GetRaceServiceInstance().DeleteRace(ctx, id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}

	render.NoContent(w, r)
}

type RaceRequest struct {
	Name string `json:"name"`
}
