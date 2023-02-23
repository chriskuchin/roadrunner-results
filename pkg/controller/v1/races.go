package v1

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

// GET /races/{id}
type racesResources struct{}

func (rs racesResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Options("/", Cors)
	r.Get("/", listRaces)
	r.Post("/", createRace)

	r.Route("/import", func(r chi.Router) {
		r.Use(handleOAuth2Creds)
		r.Get("/", importRaceAndResults)
	})

	r.Route("/{raceID}", func(r chi.Router) {
		r.Use(RaceCtx)
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

func handleOAuth2Creds(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken, err := r.Cookie("oauth-token")
		if err != nil {
			log.Error().Err(err).Send()
			url := getOAuth2Config().AuthCodeURL(r.URL.Path, oauth2.AccessTypeOffline)

			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		token, err := base64.StdEncoding.DecodeString(encodedToken.Value)
		if err != nil {
			log.Error().Err(err).Send()
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tok := &oauth2.Token{}
		err = json.Unmarshal([]byte(token), tok)
		if err != nil {
			log.Error().Err(err).Send()
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), util.OAuthTokenID, tok)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func importRaceAndResults(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test import"))
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

func RaceCtx(next http.Handler) http.Handler {
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
