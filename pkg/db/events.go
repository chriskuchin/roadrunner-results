package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type EventDao struct {
	db *sqlx.DB
}

func NewEventDAO(db *sqlx.DB) *EventDao {
	return &EventDao{
		db: db,
	}
}

func (e *EventDao) AddEvent(ctx context.Context, raceID, description string, distance int) error {
	_, err := e.db.Exec("insert into events (race_id, event_description, distance) VALUES(?, ?, ?)", raceID, description, distance)
	return err
}

func (e *EventDao) GetRaceEvents(ctx context.Context) {

}

func (e *EventDao) DeleteEvent(ctx context.Context) {

}

func (e *EventDao) GetRaceEvent(ctx context.Context) {

}
