package controller

import (
	"context"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	v1 "github.com/chriskuchin/roadrunner-results/pkg/controller/v1"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var authClient *auth.Client

func Routes(app *firebase.App, db *sqlx.DB, debug bool) chi.Router {
	var err error
	authClient, err = app.Auth(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("auth client failure")
	}

	r := chi.NewRouter()

	r.Use(authenticationMiddleware)
	if debug {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Api-Token"},
			MaxAge:         300,
		}))
	}

	r.Mount("/v1", v1.Routes(db, app, debug))
	return r
}

func authenticationMiddleware(next http.Handler) http.Handler {
	token, token_present := os.LookupEnv("API_TOKEN")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-token")
		var tokResult *auth.Token
		var tokErr error
		if apiKey != "" {
			tokResult, tokErr = authClient.VerifyIDToken(r.Context(), apiKey)
			r = r.WithContext(context.WithValue(r.Context(), util.UserToken, *tokResult))
		}

		if r.Method != http.MethodGet && r.Method != http.MethodOptions {
			if (token_present && token == apiKey) || (tokErr == nil) {
				next.ServeHTTP(w, r)
				return
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
