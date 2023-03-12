package db

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/jmoiron/sqlx"
)

type EventDao struct {
	db *sqlx.DB
}

type EventDto struct {
	RaceID      string `db:"race_id"`
	EventID     string `db:"event_id"`
	Description string `db:"event_description"`
	Distance    int    `db:"distance"`
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

	listRaceEventsQuery string = `
		select * from events
			where
				race_id = ?
	`
)

func NewEventDAO() *EventDao {
	return &EventDao{
		db: getDBInstance(),
	}
}

func (e *EventDao) AddEvent(ctx context.Context, raceID, eventID, description string, distance int) error {
	_, err := e.db.Exec(addEventQuery, raceID, eventID, description, distance)
	return err
}

func (e *EventDao) GetRaceEvents(ctx context.Context) ([]EventDto, error) {
	var results []EventDto
	return results, e.db.Select(&results, listRaceEventsQuery, util.GetRaceIDFromContext(ctx))
}

func (e *EventDao) DeleteEvent(ctx context.Context) {

}

func (e *EventDao) GetRaceEvent(ctx context.Context) {

}
