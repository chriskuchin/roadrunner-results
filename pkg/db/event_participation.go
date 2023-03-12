package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type EventParticipationDao struct {
	db *sqlx.DB
}

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

func NewEventParticipationDAO() *EventParticipationDao {
	return &EventParticipationDao{
		db: getDBInstance(),
	}
}

func (ep *EventParticipationDao) RecordEventParticipation(ctx context.Context, raceID, eventID, bibNumber string) error {
	_, err := ep.db.Exec(addEventParticipationQuery, raceID, eventID, bibNumber)
	return err
}
