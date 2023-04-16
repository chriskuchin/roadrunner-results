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

type (
	TimerRequest struct {
		Start int64 `json:"start_ts,omitempty"`
	}
)

func TimerRoutes(handler *Handler) chi.Router {
	r := chi.NewRouter()
	r.Post("/", handler.createTimer)
	r.Get("/", handler.listTimers)

	r.Route("/{timerID}", func(r chi.Router) {
		r.Use(timerCtx)
		r.Put("/", handler.startTimer)
		r.Get("/", handler.getTimer)
		r.Delete("/", handler.deleteTimer)
	})

	return r
}

func (api *Handler) createTimer(w http.ResponseWriter, r *http.Request) {
	// create and then start
	timerID, err := services.CreateTimer(r.Context(), api.db)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": fmt.Sprintf("%s", err),
		})
		return
	}

	timerR := TimerRequest{}
	err = render.DecodeJSON(r.Body, timerR)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": fmt.Sprintf("%s", err),
		})
		return
	}

	if timerR.Start != 0 {
		err = services.StartTimer(context.WithValue(r.Context(), util.TimerID, timerID), api.db, timerR.Start)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": fmt.Sprintf("%s", err),
			})
			return
		}
		render.Status(r, http.StatusAccepted)
		render.JSON(w, r, fmt.Sprintf("%d", timerR.Start))
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (api *Handler) startTimer(w http.ResponseWriter, r *http.Request) {
	timerR := TimerRequest{}
	err := render.DecodeJSON(r.Body, timerR)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": fmt.Sprintf("%s", err),
		})
		return

	}

	err = services.StartTimer(r.Context(), api.db, timerR.Start)
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
