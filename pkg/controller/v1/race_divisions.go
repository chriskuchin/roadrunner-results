package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type DivisionDefinitionPayload struct {
	Description string   `json:"display"`
	Filters     []Filter `json:"filters"`
}

type Filter struct {
	Key    string   `json:"key"`
	Values []string `json:"values,omitempty"`
}

func RaceDivisionsRoutes(handler *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.listDivisions)
	r.Post("/", handler.createDivision)

	return r
}

func (api *Handler) listDivisions(w http.ResponseWriter, r *http.Request) {
	results, err := services.ListRaceDivisions(r.Context(), api.db)
	if err != nil {
		handleBadRequest(err, w, r)
		return
	}

	render.JSON(w, r, results)
}

func (api *Handler) createDivision(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(err, w, r)
		return
	}

	var payload DivisionDefinitionPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		handleBadRequest(err, w, r)
		return
	}

	var filters []services.Filter = []services.Filter{}
	for _, v := range payload.Filters {
		filters = append(filters, services.Filter{
			Key:    v.Key,
			Values: v.Values,
		})
	}
	services.CreateDivision(ctx, api.db, payload.Description, filters)

}
