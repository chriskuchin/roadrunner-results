package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type timerResources struct{}

func (tr timerResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", startTimer)
	r.Get("/", listTimers)

	r.Route("/{timerID}", func(r chi.Router) {
		r.Use(timerCtx)
		r.Get("/", getTimer)
		r.Delete("/", deleteTimer)
	})

	return r
}

func startTimer(w http.ResponseWriter, r *http.Request) {
	services.GetTimerServiceInstance().StartTimer(r.Context(), time.Now().UnixMilli())

	w.WriteHeader(http.StatusAccepted)
}

func listTimers(w http.ResponseWriter, r *http.Request) {
	result, err := services.GetTimerServiceInstance().ListTimers(r.Context(), 10, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, result)
}

func getTimer(w http.ResponseWriter, r *http.Request) {
	result, err := services.GetTimerServiceInstance().GetTimer(r.Context())
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

func deleteTimer(w http.ResponseWriter, r *http.Request) {
	err := services.GetTimerServiceInstance().DeleteTimer(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func timerCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timerID := chi.URLParam(r, "timerID")
		if timerID == "" {
			log.Error().Msg("Unable to locate timerID")
		}

		ctx := context.WithValue(r.Context(), util.TimerID, timerID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
