package middleware

import (
	"context"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

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

func UserAuthMiddleware(db *sqlx.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			raceID := util.GetRaceIDFromContext(r.Context())
			uid := util.GetCurrentUserID(r.Context())
			ownerID, err := services.GetRaceOwnerID(r.Context(), db, raceID)
			if err != nil {
				log.Error().Err(err).Send()
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			if uid == ownerID {
				next.ServeHTTP(w, r)
				return
			}

			authorizedUsers, err := services.GetRaceAuthorizedUsers(r.Context(), db, raceID)
			if err != nil {
				log.Error().Err(err).Send()
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			if authorizedUsers[uid] {
				next.ServeHTTP(w, r)
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
		})
	}
}
