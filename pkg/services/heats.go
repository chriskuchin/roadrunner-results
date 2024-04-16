package services

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
)

type (
	AssignmentPayload struct {
		Assignments []LaneAssignment `json:"assignments"`
	}
	LaneAssignment struct {
		Lane int    `json:"lane"`
		Bib  string `json:"bib"`
	}

	heatRow struct {
		RaceID      string         `db:"race_id"`
		EventID     string         `db:"event_id"`
		TimerID     string         `db:"timer_id"`
		StartTS     int            `db:"start_ts"`
		Assignments sql.NullString `db:"assignments"`
	}

	Heat struct {
		RaceID      string
		EventID     string
		TimerID     string
		StartTS     int
		Assignments *AssignmentPayload
	}
)

const (
	// const INSERT_RACE_AUTHORIZATION = `INSERT into race_authorization (user_id, race_id, permissions)
	// VALUES(?, ?, json(?))`

	insertAssignmentsQuery = `
		INSERT INTO lane_assignments
			(race_id, event_id, timer_id, assignments)
			VALUES(?, ?, ?, json(?))
	`

	updateAssignmentsQuery = `
	UPDATE lane_assignments
		SET assignments = json(?)
		WHERE
			race_id = ? AND event_id = ? AND timer_id = ?
	`

	listHeatsQuery string = `
		select * from timers
			JOIN lane_assignments
				USING(race_id, event_id, timer_id)
		WHERE race_id = ? AND event_id = ?
	`

	DeleteLaneAssignmentsQuery string = `
	DELETE FROM lane_assignments WHERE race_id = ? AND event_id = ? AND timer_id = ?
	`
)

func DeleteHeat(ctx context.Context, db db.DB, race_id, event_id, timer_id string) error {
	_, err := db.ExecContext(ctx, DeleteLaneAssignmentsQuery, race_id, event_id, timer_id)
	if err != nil {
		return err
	}

	return DeleteTimer(ctx, db, race_id, event_id, timer_id)
}

func CreateLaneAssignment(ctx context.Context, db db.DB, race_id, event_id, timer_id string, assignments AssignmentPayload) error {
	payload, err := json.Marshal(assignments)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, insertAssignmentsQuery, race_id, event_id, timer_id, string(payload))
	return err
}

func UpdateLaneAssignments(ctx context.Context, db db.DB, race_id, event_id, timer_id string, assigments AssignmentPayload) error {
	payload, err := json.Marshal(assigments)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, updateAssignmentsQuery, payload, race_id, event_id, timer_id)

	return err
}

func ListHeatDefinitions(ctx context.Context, db db.DB, race_id, event_id string) ([]Heat, error) {
	var rows []heatRow
	err := db.SelectContext(ctx, &rows, listHeatsQuery, race_id, event_id)
	if err != nil {
		return nil, err
	}

	var result []Heat = []Heat{}
	for _, row := range rows {
		var assignments *AssignmentPayload = nil

		if row.Assignments.Valid {
			assignments = &AssignmentPayload{}
			err := json.Unmarshal([]byte(row.Assignments.String), &assignments)
			if err != nil {
				return nil, err
			}
		}

		result = append(result, Heat{
			RaceID:      row.RaceID,
			EventID:     row.EventID,
			TimerID:     row.TimerID,
			StartTS:     row.StartTS,
			Assignments: assignments,
		})
	}

	return result, nil
}
