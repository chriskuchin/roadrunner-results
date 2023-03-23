package db

import (
	"context"
	"errors"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type TimerDAO struct {
	db *sqlx.DB
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

func NewTimerDAO() *TimerDAO {
	return &TimerDAO{
		db: getDBInstance(),
	}
}

func (td *TimerDAO) StartTimer(ctx context.Context, timerID string, start int64) error {
	_, err := td.db.Exec(startTimer, timerID, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), start)
	return err
}

func (td *TimerDAO) ListTimers(ctx context.Context, limit, offset int) ([]TimerRow, error) {
	var results []TimerRow
	err := td.db.Select(&results, listEventTimersQuery, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), limit, offset)

	return results, err
}

func (td *TimerDAO) GetTimer(ctx context.Context) (*TimerRow, error) {
	var result []*TimerRow
	err := td.db.Select(&result, getEventTimerQuery, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), util.GetTimerIDFromContext(ctx))

	if len(result) == 0 {
		return nil, errors.New("timer not found")
	}

	return result[0], err
}

func (td *TimerDAO) GetActiveTimerStart(ctx context.Context) (int64, string, error) {
	row := td.db.QueryRow(getActiveEventTimerQuery, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx))
	if row.Err() != nil {
		return 0, "", row.Err()
	}

	var start int64
	var timerID string
	row.Scan(&timerID, &start)
	return start, timerID, nil
}

func (td *TimerDAO) DeleteTimer(ctx context.Context) error {
	_, err := td.db.Exec(deleteEventTimerQuery, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), util.GetTimerIDFromContext(ctx))
	return err
}

func (td *TimerDAO) RecordTime(ctx context.Context, timerID string, result int64) error {
	_, err := td.db.Exec(appendTimerResult, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), timerID, result)
	log.Error().Err(err).Send()
	return err
}
