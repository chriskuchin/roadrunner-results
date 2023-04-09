package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type (
	EventObject struct {
		EventID     string `json:"eventId"`
		Description string `json:"description"`
		Distance    int    `json:"distance"`
	}

	EventRow struct {
		RaceID      string `db:"race_id"`
		EventID     string `db:"event_id"`
		Description string `db:"event_description"`
		Distance    int    `db:"distance"`
	}
)

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

func AddEvent(ctx context.Context, db *sqlx.DB, raceID, description string, distance int) (string, error) {
	id := uuid.NewString()
	_, err := db.Exec(addEventQuery, raceID, id, description, distance)
	return id, err
}

func GetRaceEvents(ctx context.Context, db *sqlx.DB) ([]EventObject, error) {
	var dbEvents []EventRow
	err := db.Select(&dbEvents, listRaceEventsQuery, util.GetRaceIDFromContext(ctx))
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	results := []EventObject{}
	for _, event := range dbEvents {
		results = append(results, EventObject{
			EventID:     event.EventID,
			Description: event.Description,
			Distance:    event.Distance,
		})
	}

	return results, nil
}

func GetRaceEvent(ctx context.Context) {

}

func DeleteRaceEvent(ctx context.Context) {

}

func GetRaceEventCount(ctx context.Context, db *sqlx.DB) (int, error) {
	row := db.QueryRow(selectRaceEventCountQuery, util.GetRaceIDFromContext(ctx))
	if row.Err() != nil {
		log.Error().Msg("err")
		return 0, row.Err()
	}

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}