package middleware

import (
	"context"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func EventCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventID := chi.URLParam(r, "eventID")
		if eventID == "" {
			log.Error().Msg("Unable to locate eventID")
		}

		ctx := context.WithValue(r.Context(), util.EventID, eventID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
