package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/google"
	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func RacesRoutes(handler *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.listRaces)
	r.Post("/", handler.createRace)

	r.Route("/import", func(r chi.Router) {
		r.Use(google.HandleOAuth2Creds)
		r.Get("/", handler.importRaceAndResults)
	})

	r.Route("/{raceID}", func(r chi.Router) {
		r.Use(raceCtx)
		r.Use(handler.userIsAuthorized)

		r.Delete("/", handler.deleteRace)
		r.Get("/", handler.getRace)

		r.Mount("/volunteers", VolunteerRoutes(handler))
		r.Mount("/participants", ParticipantsRoutes(handler))
		r.Mount("/events", EventsRoutes(handler))
	})

	return r
}

func (api *Handler) importRaceAndResults(w http.ResponseWriter, r *http.Request) {
	sheetId := r.URL.Query().Get("sheetId")
	if sheetId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no sheetId provided"))
		return
	}
	services.ImportFromSheet(r.Context(), api.db, sheetId)

	w.Write([]byte("test import, " + sheetId))
}

func (api *Handler) getRace(w http.ResponseWriter, r *http.Request) {
	race, err := services.GetRace(r.Context(), api.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, race)
}

// GET /races
func (api *Handler) listRaces(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	races, _ := services.ListRaces(ctx, api.db)

	render.JSON(w, r, races)
}

// POST /races
func (api *Handler) createRace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	raceRequest := &RaceRequest{}
	if err := render.DecodeJSON(r.Body, raceRequest); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": fmt.Sprintf("%v", err),
		})
	}

	raceID, err := services.CreateRace(ctx, api.db, raceRequest.Name)
	if err != nil {
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
func (api *Handler) deleteRace(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "raceID")
	ctx := r.Context()
	err := services.DeleteRace(ctx, api.db, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete race")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}

	render.NoContent(w, r)
}

func (api *Handler) userIsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		raceID := util.GetRaceIDFromContext(r.Context())
		uid := util.GetCurrentUserID(r.Context())
		ownerID, err := services.GetRaceOwnerID(r.Context(), api.db, raceID)
		if err != nil {
			log.Error().Err(err).Send()
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		if uid == ownerID {
			next.ServeHTTP(w, r)
			return
		}

		log.Debug().Str("uid", uid).Str("owner", ownerID).Send()
		w.WriteHeader(http.StatusUnauthorized)
	})
}

func raceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raceID := chi.URLParam(r, "raceID")
		if raceID == "" {
			log.Error().Msg("Unable to locate RaceID")
		}

		ctx := context.WithValue(r.Context(), util.RaceID, raceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type RaceRequest struct {
	Name string `json:"name"`
}
