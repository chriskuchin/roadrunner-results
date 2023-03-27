package services

import (
	"context"
	"errors"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type TimerResult struct {
	TimerID string `json:"id"`
	RaceID  string `json:"race_id"`
	EventID string `json:"event_id"`
	Start   int64  `json:"timer_start"`
}

func StartTimer(ctx context.Context, db *sqlx.DB, start int64) {
	timerID := uuid.New().String()
	db.Exec(startTimer, timerID, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), start)
}

func ListTimers(ctx context.Context, db *sqlx.DB, limit, offset int) ([]TimerResult, error) {
	var rows []TimerRow
	err := db.Select(&rows, listEventTimersQuery, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), limit, offset)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	var results []TimerResult
	for _, row := range rows {

		results = append(results, TimerResult{
			TimerID: row.TimerID,
			RaceID:  row.RaceID,
			EventID: row.EventID,
			Start:   row.Start,
		})
	}

	return results, nil
}

func DeleteTimer(ctx context.Context, db *sqlx.DB) error {
	_, err := db.Exec(deleteEventTimerQuery, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), util.GetTimerIDFromContext(ctx))
	return err

}

func GetTimer(ctx context.Context, db *sqlx.DB) (*TimerResult, error) {
	var result []*TimerRow
	err := db.Select(&result, getEventTimerQuery, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), util.GetTimerIDFromContext(ctx))
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

type TimerRow struct {
	TimerID string `db:"timer_id"`
	RaceID  string `db:"race_id"`
	EventID string `db:"event_id"`
	Start   int64  `db:"start_ts"`
}

const (
	startTimer string = `
		insert into timers (
			timer_id,
			race_id,
			event_id,
			start_ts
		)
		VALUES(?, ?, ?, ?)
	`

	listEventTimersQuery string = `
		select * from timers
		where
			race_id = ?
			AND
			event_id = ?
		limit ? offset ?
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
		select timer_id, max(start_ts) from timers
		Where
			race_id = ?
			AND
			event_id = ?
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
)

func GetActiveTimerStart(ctx context.Context, db *sqlx.DB) (int64, string, error) {
	row := db.QueryRow(getActiveEventTimerQuery, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx))
	if row.Err() != nil {
		return 0, "", row.Err()
	}

	var start int64
	var timerID string
	row.Scan(&timerID, &start)
	return start, timerID, nil
}

func RecordTime(ctx context.Context, db *sqlx.DB, timerID string, result int64) error {
	_, err := db.Exec(appendTimerResult, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), timerID, result)
	return err
}
