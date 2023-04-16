package controller

import (
	v1 "github.com/chriskuchin/roadrunner-results/pkg/controller/v1"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
)

func Routes(db *sqlx.DB, debug bool) chi.Router {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	}))

	r.Mount("/v1", v1.Routes(db, debug))

	return r
}
