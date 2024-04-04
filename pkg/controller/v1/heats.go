package v1

import (
	"net/http"

	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

func HandleHeatsCreate(db *sqlx.DB) http.HandlerFunc {
	type assignment struct {
		Lane int    `json:"lane"`
		Bib  string `json:"bib"`
	}
	type request struct {
		Assignments []assignment `json:"assignments"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		raceID := util.GetRaceIDFromContext(ctx)
		eventID := util.GetEventIDFromContext(ctx)

		var reqPayload request = request{}
		err := render.DecodeJSON(r.Body, &reqPayload)
		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		timerID, err := services.CreateTimer(ctx, db, raceID, eventID)
		if err != nil {
			apiutil.HandleServerError(err, w, r)
			return
		}

		var laneAssignments services.AssignmentPayload = services.AssignmentPayload{
			Assignments: []services.LaneAssignment{},
		}

		if reqPayload.Assignments != nil {
			for _, assignment := range reqPayload.Assignments {
				laneAssignments.Assignments = append(laneAssignments.Assignments, services.LaneAssignment(assignment))
			}
		}

		err = services.CreateLaneAssignment(ctx, db, raceID, eventID, timerID, laneAssignments)
		if err != nil {
			apiutil.HandleServerError(err, w, r)
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, map[string]string{
			"timer_id": timerID,
		})
	}
}

func HandleHeatUpdate(db *sqlx.DB) http.HandlerFunc {
	type assignment struct {
		Lane int    `json:"lane"`
		Bib  string `json:"bib"`
	}
	type request struct {
		Assignments []assignment `json:"assignments"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req = request{}
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		var assignmentsPayload = services.AssignmentPayload{
			Assignments: []services.LaneAssignment{},
		}

		for _, assgnmt := range req.Assignments {
			assignmentsPayload.Assignments = append(assignmentsPayload.Assignments, services.LaneAssignment{
				Lane: assgnmt.Lane,
				Bib:  assgnmt.Bib,
			})
		}

		err = services.UpdateLaneAssignments(ctx, db, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), util.GetTimerIDFromContext(ctx), assignmentsPayload)
		if err != nil {
			apiutil.HandleServerError(err, w, r)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func HandleHeatDelete(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func HandleHeatsList(db *sqlx.DB) http.HandlerFunc {
	type assignments struct {
		Lane int    `json:"lane"`
		Bib  string `json:"bib"`
	}
	type response struct {
		RaceID      string        `json:"race_id"`
		EventID     string        `json:"event_id"`
		TimerID     string        `json:"timer_id"`
		Start       int           `json:"start"`
		Assignments []assignments `json:"assignments"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		race_id := util.GetRaceIDFromContext(ctx)
		event_id := util.GetEventIDFromContext(ctx)

		result, err := services.ListHeatDefinitions(ctx, db, race_id, event_id)
		if err != nil {
			apiutil.HandleServerError(err, w, r)
			return
		}

		var heats []response = []response{}
		for _, heat := range result {

			var lanes []assignments = nil
			if heat.Assignments != nil {
				lanes = []assignments{}
				for _, lane := range heat.Assignments.Assignments {
					lanes = append(lanes, assignments{
						Lane: lane.Lane,
						Bib:  lane.Bib,
					})

				}
			}

			heats = append(heats, response{
				RaceID:      heat.RaceID,
				EventID:     heat.EventID,
				TimerID:     heat.TimerID,
				Start:       heat.StartTS,
				Assignments: lanes,
			})
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, heats)
	}
}
