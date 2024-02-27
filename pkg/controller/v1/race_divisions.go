package v1

import (
	"encoding/json"
	"io"
	"net/http"

	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

func HandleDivisionsList(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		results, err := services.ListRaceDivisions(ctx, db, util.GetRaceIDFromContext(ctx))
		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		render.JSON(w, r, results)
	}
}

func HandleDivisionsCreate(db *sqlx.DB) http.HandlerFunc {
	type Filter struct {
		Key    string   `json:"key"`
		Values []string `json:"values,omitempty"`
	}

	type DivisionDefinitionPayload struct {
		Description string   `json:"display"`
		Filters     []Filter `json:"filters"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		var payload DivisionDefinitionPayload
		err = json.Unmarshal(body, &payload)
		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		var filters []services.Filter = []services.Filter{}
		for _, v := range payload.Filters {
			filters = append(filters, services.Filter{
				Key:    v.Key,
				Values: v.Values,
			})
		}
		services.CreateDivision(ctx, db, util.GetRaceIDFromContext(ctx), payload.Description, filters)
	}
}
