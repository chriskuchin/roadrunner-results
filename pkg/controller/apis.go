package controller

import (
	"net/http"
	"os"

	v1 "github.com/chriskuchin/roadrunner-results/pkg/controller/v1"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func Routes(db *sqlx.DB, debug bool) chi.Router {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Api-Token"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(auth)

	r.Mount("/v1", v1.Routes(db, debug))

	return r
}

func auth(next http.Handler) http.Handler {
	token, token_present := os.LookupEnv("API_TOKEN")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodOptions {
			apiKey := r.Header.Get("x-api-token")
			log.Info().Str("ApiKey", apiKey).Str("token", token).Bool("present", token_present).Send()
			if token_present && token == apiKey {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
