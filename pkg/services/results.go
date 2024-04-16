package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
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

func InsertResults(ctx context.Context, db *db.DBLayer, raceID, eventID, bibNumber, result string) error {
	_, err := db.ExecContext(ctx, addResultQuery, raceID, eventID, bibNumber, result)
	return err
}

func InsertPartialResult(ctx context.Context, db *db.DBLayer, raceID, eventID string, result int64, timerID string) error {
	var err error
	if timerID == "" {
		_, err = db.ExecContext(ctx, insertPartialQuery, raceID, eventID, nil, result)
	} else {
		_, err = db.ExecContext(ctx, insertPartialQuery, raceID, eventID, timerID, result)
	}
	return err
}
