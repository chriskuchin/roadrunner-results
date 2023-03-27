package v1

import (
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	db    *sqlx.DB
	debug bool
}

func Routes(db *sqlx.DB, debug bool) chi.Router {
	r := chi.NewRouter()
	handler := &Handler{
		db:    db,
		debug: debug,
	}
	r.Mount("/races", RacesRoutes(handler))
	r.Mount("/google", GoogleRoutes(handler))
	return r
}
