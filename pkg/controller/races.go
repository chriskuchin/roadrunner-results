package controller

import "github.com/go-chi/chi"

type RacesResources struct{}

func (rs RacesResources) Routes() chi.Router {
	r := chi.NewRouter()

	return r
}
