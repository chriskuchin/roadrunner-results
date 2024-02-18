package middleware

import (
	"context"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func TimerCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timerID := chi.URLParam(r, "timerID")
		if timerID == "" {
			log.Error().Msg("Unable to locate timerID")
		}

		ctx := context.WithValue(r.Context(), util.TimerID, timerID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
