package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HealthcheckResources struct{}

func (rs HealthcheckResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", getHealthcheck)

	return r
}

func getHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
