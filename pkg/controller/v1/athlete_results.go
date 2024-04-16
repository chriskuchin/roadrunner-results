package v1

import (
	"fmt"
	"net/http"

	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/go-chi/render"
)

func HandleAthleteResultsSearch(db *db.DBLayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var filter map[string]string = map[string]string{}

		if r.URL.Query().Has("gender") {
			filter["gender"] = r.URL.Query().Get("gender")
		}

		if r.URL.Query().Has("first_name") {
			filter["first_name"] = r.URL.Query().Get("first_name")
		}

		if r.URL.Query().Has("last_name") {
			filter["last_name"] = r.URL.Query().Get("last_name")
		}

		if r.URL.Query().Has("birth_year") {
			filter["birth_year"] = r.URL.Query().Get("birth_year")
		}

		if r.URL.Query().Has("type") {
			filter["type"] = r.URL.Query().Get("type")
		}

		if len(filter) < 2 {
			apiutil.HandleBadRequest(fmt.Errorf("not enough filters"), w, r)
			return
		}

		results, err := services.FindAthleteResults(r.Context(), db, filter)
		if err != nil {
			apiutil.HandleServerError(err, w, r)
			return
		}

		render.JSON(w, r, results)
	}
}
