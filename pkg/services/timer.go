package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type (
	TimerRow struct {
		Results     int            `db:"count"`
		TimerID     string         `db:"timer_id"`
		RaceID      string         `db:"race_id"`
		EventID     string         `db:"event_id"`
		Start       int64          `db:"start_ts"`
		Assignments sql.NullString `db:"assignments"`
	}

	TimerResult struct {
		Results     int                 `json:"count"`
		TimerID     string              `json:"id"`
		RaceID      string              `json:"race_id"`
		EventID     string              `json:"event_id"`
		Start       int64               `json:"timer_start"`
		Assignments []AssignmentsResult `json:"assignments,omitempty"`
	}

	AssignmentsResult struct {
		Lane uint   `json:"lane"`
		Bib  string `json:"bib"`
	}
)

const (
	startTimerQuery string = `
		update timers set start_ts = ?
		WHERE
			timer_id = ?
				AND
			event_id = ?
				AND
			race_id = ?
	`

	listEventTimersQuery string = `
		select
			t.*, json_extract(la.assignments, "$") as assignments, count(distinct r.result) as count from timers as t
		left outer join results r
			using(timer_id, event_id, race_id)
		left outer join lane_assignments la
			using(race_id, event_id, timer_id)
		where
			t.race_id = ?
			AND
			t.event_id = ?
		group by t.timer_id, t.event_id, t.race_id
		order by t.rowid ASC
	`

	getEventTimerQuery string = `
		select * from timers
		Where
			race_id = ?
			AND
			event_id = ?
			AND
			timer_id = ?
	`

	getActiveEventTimerQuery string = `
		select timer_id, start_ts from timers
		Where
			race_id = ?
			AND
			event_id = ?
		ORDER BY start_ts DESC
		Limit 1
	`

	deleteEventTimerQuery string = `
		delete from timers
		Where
			race_id = ?
			AND
			event_id = ?
			AND
			timer_id = ?
	`

	appendTimerResult string = `
	insert into timer_results (
		race_id,
		event_id,
		timer_id,
		result
	) VALUES (?, ?, ?, ?)
	`

	createTimerQuery string = `
	insert into timers (
		timer_id,
		event_id,
		race_id,
		start_ts
	) VALUES(?, ?, ?, 0)`
)

func CreateTimer(ctx context.Context, db *sqlx.DB, raceID, eventID string) (string, error) {
	timerID := uuid.New().String()
	_, err := db.Exec(createTimerQuery, timerID, eventID, raceID)
	return timerID, err
}

func StartTimer(ctx context.Context, db *sqlx.DB, raceID, eventID, timerID string, start int64) error {
	res, err := db.Exec(startTimerQuery, start, timerID, eventID, raceID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("failed to start timer")
	}

	return nil
}

func ListTimers(ctx context.Context, db *sqlx.DB, raceID, eventID string, limit, offset int) ([]TimerResult, error) {
	var rows []TimerRow
	err := db.Select(&rows, listEventTimersQuery, raceID, eventID, limit, offset)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	var results []TimerResult
	for _, row := range rows {
		resolvedRow := TimerResult{
			Results: row.Results,
			TimerID: row.TimerID,
			EventID: row.EventID,
			RaceID:  row.RaceID,
			Start:   row.Start,
		}

		var assignments map[string][]AssignmentsResult
		if row.Assignments.Valid {
			err := json.Unmarshal([]byte(row.Assignments.String), &assignments)
			if err != nil {
				log.Warn().Err(err).Msg("failed deserializing lane assignments")
			} else {
				resolvedRow.Assignments = assignments["assignments"]
			}
		}

		results = append(results, resolvedRow)
	}

	return results, nil
}

func DeleteTimer(ctx context.Context, db *sqlx.DB, raceID, eventID, timerID string) error {
	_, err := db.Exec(deleteEventTimerQuery, raceID, eventID, timerID)
	return err

}

func GetTimer(ctx context.Context, db *sqlx.DB, raceID, eventID, timerID string) (*TimerResult, error) {
	var result []*TimerRow
	err := db.Select(&result, getEventTimerQuery, raceID, eventID, timerID)
	if err != nil {
		return nil, err
	} else if len(result) == 0 {
		return nil, errors.New("timer not found")
	}

	return &TimerResult{
		RaceID:  result[0].RaceID,
		EventID: result[0].EventID,
		TimerID: result[0].TimerID,
		Start:   result[0].Start,
	}, nil
}

func GetActiveTimerStart(ctx context.Context, db *sqlx.DB, raceID, eventID string) (int64, string, error) {
	row := db.QueryRow(getActiveEventTimerQuery, raceID, eventID)
	if row.Err() != nil {
		return 0, "", row.Err()
	}

	var start int64
	var timerID string
	row.Scan(&timerID, &start)
	return start, timerID, nil
}

func RecordTime(ctx context.Context, db *sqlx.DB, raceID, eventID, timerID string, result int64) error {
	_, err := db.Exec(appendTimerResult, raceID, eventID, timerID, result)
	return err
}
