package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type (
	TimerRow struct {
		Results int    `db:"count"`
		TimerID string `db:"timer_id"`
		RaceID  string `db:"race_id"`
		EventID string `db:"event_id"`
		Start   int64  `db:"start_ts"`
	}

	TimerResult struct {
		Results int    `json:"count"`
		TimerID string `json:"id"`
		RaceID  string `json:"race_id"`
		EventID string `json:"event_id"`
		Start   int64  `json:"timer_start"`
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
			t.*, count(distinct r.result) as count from timers as t
		left outer join results r
			using(timer_id, event_id, race_id)
		where
			t.race_id = ?
			AND
			t.event_id = ?
		group by t.timer_id, t.event_id, t.race_id
		order by t.start_ts ASC
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

func CreateTimer(ctx context.Context, db *sqlx.DB) (string, error) {
	timerID := uuid.New().String()
	_, err := db.Exec(createTimerQuery, timerID, util.GetEventIDFromContext(ctx), util.GetRaceIDFromContext(ctx))
	return timerID, err
}

func StartTimer(ctx context.Context, db *sqlx.DB, start int64) error {
	res, err := db.Exec(startTimerQuery, start, util.GetTimerIDFromContext(ctx), util.GetEventIDFromContext(ctx), util.GetRaceIDFromContext(ctx))
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

func ListTimers(ctx context.Context, db *sqlx.DB, limit, offset int) ([]TimerResult, error) {
	var rows []TimerRow
	err := db.Select(&rows, listEventTimersQuery, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), limit, offset)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	var results []TimerResult
	for _, row := range rows {

		results = append(results, TimerResult(row))
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
