package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func UsersRoutes(handler *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.getCurrentUser)
	r.Post("/", handler.registerUser)

	return r
}

func (api *Handler) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (api *Handler) registerUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
