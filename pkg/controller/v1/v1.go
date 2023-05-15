package v1

import (
	"net/http"

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
	r.Use(auth)
	r.Mount("/races", RacesRoutes(handler))
	r.Mount("/google", GoogleRoutes(handler))
	return r
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
