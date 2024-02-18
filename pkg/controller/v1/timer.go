package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func HandleTimersCreate(db *sqlx.DB) http.HandlerFunc {
	type TimerRequest struct {
		Start int64 `json:"start_ts,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// create and then start
		timerID, err := services.CreateTimer(r.Context(), db)
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
			err = services.StartTimer(context.WithValue(r.Context(), util.TimerID, timerID), db, timerR.Start)
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

func HandleTimerStart(db *sqlx.DB) http.HandlerFunc {
	type TimerRequest struct {
		Start int64 `json:"start_ts,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		timerR := TimerRequest{}
		err := render.DecodeJSON(r.Body, timerR)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": fmt.Sprintf("%s", err),
			})
			return

		}

		err = services.StartTimer(r.Context(), db, timerR.Start)
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

func HandleTimersList(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := services.ListTimers(r.Context(), db, 10, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, result)
	}
}

func HandleTimerGet(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := services.GetTimer(r.Context(), db)
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

func HandleTimerDelete(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := services.DeleteTimer(r.Context(), db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
