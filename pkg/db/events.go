package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type EventDao struct {
	db *sqlx.DB
}

const (
	addEventQuery string = `
		insert into events (
			race_id,
			event_id,
			event_description,
			distance
			)
			VALUES(?, ?, ?, ?)
	`
)

func NewEventDAO(db *sqlx.DB) *EventDao {
	return &EventDao{
		db: db,
	}
}

func (e *EventDao) AddEvent(ctx context.Context, raceID, eventID, description string, distance int) error {
	_, err := e.db.Exec(addEventQuery, raceID, eventID, description, distance)
	return err
}

func (e *EventDao) GetRaceEvents(ctx context.Context) {

}

func (e *EventDao) DeleteEvent(ctx context.Context) {

}

func (e *EventDao) GetRaceEvent(ctx context.Context) {

}
