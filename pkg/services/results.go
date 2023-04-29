package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/jmoiron/sqlx"
)

const (
	addResultQuery = `
		insert into results (
			race_id,
			event_id,
			bib_number,
			result
		)
		VALUES (?, ?, ?, ?)
	`

	insertPartialQuery string = `
		insert into results (
			race_id,
			event_id,
			timer_id,
			bib_number,
			result
		)
		VALUES(?, ?, ?, NULL, ?)
	`
)

func InsertResults(ctx context.Context, db *sqlx.DB, raceID, eventID, bibNumber, result string) error {
	_, err := db.Exec(addResultQuery, raceID, eventID, bibNumber, result)
	return err
}

func InsertPartialResult(ctx context.Context, result int64, timerID string) error {
	var err error
	if timerID == "" {
		_, err = util.GetDB(ctx).Exec(insertPartialQuery, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), nil, result)
	} else {
		_, err = util.GetDB(ctx).Exec(insertPartialQuery, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), timerID, result)
	}
	return err
}
