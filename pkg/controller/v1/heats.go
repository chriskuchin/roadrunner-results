package v1

import (
	"context"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func HeatsRoutes(handler *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", listHeats)

	r.Route("/{heatId}", func(r chi.Router) {
		r.Use(heatsCtx)

		r.Get("/", getHeat)
		r.Put("/", associateHeatWithTimer)
		r.Post("/", createHeat)
	})

	return r
}

func listHeats(w http.ResponseWriter, r *http.Request) {

}

func getHeat(w http.ResponseWriter, r *http.Request) {

}

func createHeat(w http.ResponseWriter, r *http.Request) {

}

func associateHeatWithTimer(w http.ResponseWriter, r *http.Request) {

}

func heatsCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		heatID := chi.URLParam(r, "heatId")
		if heatID == "" {
			log.Error().Msg("Unable to locate heatId")
		}

		ctx := context.WithValue(r.Context(), util.HeatID, heatID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
