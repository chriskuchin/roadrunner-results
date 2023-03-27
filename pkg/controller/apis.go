package controller

import (
	v1 "github.com/chriskuchin/roadrunner-results/pkg/controller/v1"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

func Routes(db *sqlx.DB, debug bool) chi.Router {
	r := chi.NewRouter()
	r.Mount("/v1", v1.Routes(db, debug))

	return r
}
