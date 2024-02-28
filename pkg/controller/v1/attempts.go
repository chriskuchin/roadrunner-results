package v1

import (
	"net/http"

	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

func HandleEventAttemptsCreate(db *sqlx.DB) http.HandlerFunc {
	type request struct {
		BibNumber     string  `json:"bib"`
		Distance      float32 `json:"distance"`
		AttemptNumber int     `json:"attempt_number"`
	}
	type attempt struct {
		AttemptNumber int     `json:"attempt_number"`
		Distance      float32 `json:"distance"`
	}
	type response struct {
		Attempts []attempt `json:"attempts"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var request = &request{}
		err := render.DecodeJSON(r.Body, request)
		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		err = services.RecordAttempt(ctx, db, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), request.BibNumber, request.AttemptNumber, request.Distance)
		if err != nil {
			apiutil.HandleServerError(err, w, r)
			return
		}

		err = services.UpsertBestAttemptResults(ctx, db, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), request.BibNumber)
		if err != nil {
			apiutil.HandleServerError(err, w, r)
			return
		}

		attempts, err := services.ListAttempts(ctx, db, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), request.BibNumber)
		if err != nil {
			apiutil.HandleServerError(err, w, r)
			return
		}

		var result = response{
			Attempts: []attempt{},
		}
		for _, a := range attempts {
			result.Attempts = append(result.Attempts, attempt{
				AttemptNumber: a.Attempt,
				Distance:      a.Result,
			})
		}
		render.JSON(w, r, result)
	}
}

func HandleEventAttemptsList(db *sqlx.DB) http.HandlerFunc {
	type attemptResponse struct {
		AttemptNumber int     `json:"attempt_number"`
		Result        float32 `json:"result"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var response []attemptResponse = []attemptResponse{}
		ctx := r.Context()
		attempts, err := services.ListAttempts(ctx, db, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), util.GetBibNumberFromContext(ctx))
		if err != nil {
			apiutil.HandleServerError(err, w, r)
		}

		for _, attempt := range attempts {
			response = append(response, attemptResponse{
				AttemptNumber: attempt.Attempt,
				Result:        attempt.Result,
			})
		}

		render.JSON(w, r, response)
	}
}
