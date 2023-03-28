package services

import (
	"context"

	"github.com/jmoiron/sqlx"
)

const (
	addEventParticipationQuery string = `
		insert into event_participation(
			race_id,
			event_id,
			bib_number
		)
		values (?,?,?)
	`
)

func RecordEventParticipation(ctx context.Context, db *sqlx.DB, raceID, eventID, bibNumber string) error {
	_, err := db.Exec(addEventParticipationQuery, raceID, eventID, bibNumber)
	return err
}
