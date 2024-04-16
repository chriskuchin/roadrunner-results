package v1

import (
	"fmt"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func HandleTimersCreate(db *db.DBLayer) http.HandlerFunc {
	type assignments struct {
		Lane int    `json:"lane,omitempty"`
		Bib  string `json:"bib,omitempty"`
	}
	type TimerRequest struct {
		Start       int64         `json:"start_ts,omitempty"`
		Assignments []assignments `json:"assignments,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// create and then start
		ctx := r.Context()
		timerID, err := services.CreateTimer(ctx, db, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx))
		if err != nil {
			log.Error().Err(err).Send()
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{
				"error": fmt.Sprintf("%s", err),
			})
			return
		}

		timerR := TimerRequest{}
		err = render.DecodeJSON(r.Body, &timerR)
		if err != nil {
			log.Error().Err(err).Send()
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": fmt.Sprintf("%s", err),
			})
			return
		}

		if timerR.Start != 0 {
			err = services.StartTimer(ctx, db, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), timerID, timerR.Start)
			if err != nil {
				log.Error().Err(err).Send()
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, map[string]string{
					"error": fmt.Sprintf("%s", err),
				})
				return
			}
			result := map[string]string{
				"id": timerID,
			}
			render.Status(r, http.StatusAccepted)
			render.JSON(w, r, result)
			return
		}

		render.Status(r, http.StatusNoContent)
	}
}

func HandleTimerStart(db *db.DBLayer) http.HandlerFunc {
	type TimerRequest struct {
		Start int64 `json:"start_ts,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		timerR := TimerRequest{}
		err := render.DecodeJSON(r.Body, &timerR)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": fmt.Sprintf("%s", err),
			})
			return
		}

		err = services.StartTimer(ctx, db, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), util.GetTimerIDFromContext(ctx), timerR.Start)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{
				"error": fmt.Sprintf("%s", err),
			})
			return
		}

		render.Status(r, http.StatusAccepted)
		render.JSON(w, r, fmt.Sprintf("%d", timerR.Start))
	}
}

func HandleTimersList(db *db.DBLayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		result, err := services.ListTimers(ctx, db, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), 10, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, result)
	}
}

func HandleTimerGet(db *db.DBLayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		result, err := services.GetTimer(ctx, db, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), util.GetTimerIDFromContext(ctx))
		if err != nil {
			if err.Error() == "timer not found" {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		render.JSON(w, r, result)
	}
}

func HandleTimerDelete(db *db.DBLayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := services.DeleteTimer(ctx, db, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), util.GetTimerIDFromContext(ctx))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
