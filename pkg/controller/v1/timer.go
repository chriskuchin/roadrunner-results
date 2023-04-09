package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func TimerRoutes(handler *Handler) chi.Router {
	r := chi.NewRouter()
	r.Post("/", handler.startTimer)
	r.Get("/", handler.listTimers)

	r.Route("/{timerID}", func(r chi.Router) {
		r.Use(timerCtx)
		r.Get("/", handler.getTimer)
		r.Delete("/", handler.deleteTimer)
	})

	return r
}

func (api *Handler) startTimer(w http.ResponseWriter, r *http.Request) {
	start := time.Now().UnixMilli()
	services.StartTimer(r.Context(), api.db, start)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%d", start)))
}

func (api *Handler) listTimers(w http.ResponseWriter, r *http.Request) {
	result, err := services.ListTimers(r.Context(), api.db, 10, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, result)
}

func (api *Handler) getTimer(w http.ResponseWriter, r *http.Request) {
	result, err := services.GetTimer(r.Context(), api.db)
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

func (api *Handler) deleteTimer(w http.ResponseWriter, r *http.Request) {
	err := services.DeleteTimer(r.Context(), api.db)
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
