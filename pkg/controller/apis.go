package controller

import (
	"net/http"
	"os"
	"context"

	firebase "firebase.google.com/go/v4"
	v1 "github.com/chriskuchin/roadrunner-results/pkg/controller/v1"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
)

func Routes(app *firebase.App, db *sqlx.DB, debug bool) chi.Router {
	app.Auth(context.Background())

	r := chi.NewRouter()

	if !debug {
		r.Use(auth)
	} else {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Api-Token"},
			MaxAge:         300,
		}))
	}

	r.Mount("/v1", v1.Routes(db, debug))
	return r
}

func auth(next http.Handler) http.Handler {
	token, token_present := os.LookupEnv("API_TOKEN")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodOptions {
			apiKey := r.Header.Get("x-api-token")
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
