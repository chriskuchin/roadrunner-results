package v1

import (
	"github.com/go-chi/chi"
)

type Resources struct{}

func (rs Resources) Routes() chi.Router {
	r := chi.NewRouter()
	r.Mount("/races", racesResources{}.Routes())
	r.Mount("/google", googleResources{}.Routes())
	return r
}
