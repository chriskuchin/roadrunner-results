package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func HandleRaceImportSheet(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sheetId := r.URL.Query().Get("sheetId")
		if sheetId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no sheetId provided"))
			return
		}
		services.ImportFromSheet(r.Context(), db, sheetId)

		w.Write([]byte("test import, " + sheetId))
	}
}

func HandleRaceGet(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		race, err := services.GetRace(ctx, db, util.GetRaceIDFromContext(ctx))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, race)
	}
}

func HandleRacesList(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		limit := 20
		if r.URL.Query().Has("limit") {
			limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
			if err != nil {
				apiutil.HandleBadRequest(err, w, r)
				return
			}
		}

		offset := 0
		if r.URL.Query().Has("offset") {
			offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
			if err != nil {
				apiutil.HandleBadRequest(err, w, r)
				return
			}
		}

		races, err := services.ListRaces(r.Context(), db, limit, offset)
		if err != nil {
			apiutil.HandleServerError(err, w, r)
			return
		}

		render.JSON(w, r, races)
	}
}

func HandleRacesCreate(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		raceRequest := &RaceRequest{}
		if err := render.DecodeJSON(r.Body, raceRequest); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": fmt.Sprintf("%v", err),
			})
		}

		var date time.Time = time.Now()
		if raceRequest.Date != "" {
			var err error
			date, err = time.Parse("2006-01-02", raceRequest.Date)
			if err != nil {
				apiutil.HandleBadRequest(err, w, r)
				return
			}
		}

		if raceRequest.ImportURL != "" {
			raceID, err := services.ImportRaceFromURL(ctx, db, raceRequest.ImportURL, date, raceRequest.Name)
			if err != nil {
				apiutil.HandleServerError(err, w, r)
				return
			}

			render.Status(r, http.StatusCreated)
			render.JSON(w, r, map[string]string{
				"id": raceID,
			})

			return
		}

		raceID, err := services.CreateRace(ctx, db, raceRequest.Name, date)
		if err != nil {
			apiutil.HandleServerError(err, w, r)
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, map[string]string{
			"id": raceID,
		})
	}
}

func HandleRaceDelete(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "raceID")
		ctx := r.Context()
		err := services.DeleteRace(ctx, db, id)
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
}

type RaceRequest struct {
	Name      string `json:"name"`
	Date      string `json:"date"`
	ImportURL string `json:"url"`
}
