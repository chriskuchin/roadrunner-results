package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type (
	EventObject struct {
		EventID     string  `json:"eventId"`
		Description string  `json:"description"`
		Type        string  `json:"type"`
		Distance    float64 `json:"distance"`
	}

	EventRow struct {
		RaceID      string  `db:"race_id"`
		EventID     string  `db:"event_id"`
		Description string  `db:"event_description"`
		Type        string  `db:"event_type"`
		Distance    float64 `db:"distance"`
	}
)

const (
	addEventQuery string = `
		insert into events (
			race_id,
			event_id,
			event_description,
			event_type,
			distance
			)
			VALUES(?, ?, ?, ?, ?)
	`

	listRaceEventsQuery string = `
		select * from events
			where
				race_id = ?
	`

	selectEventQuery string = `
		select * from events
			where
				race_id = ?
				and
				event_id = ?
	`

	deleteEventQuery string = `
		delete from events
			where
				race_id = ?
				and
				event_id = ?
	`
)

func AddEvent(ctx context.Context, db *db.DBLayer, raceID, description, eventType string, distance float64) (string, error) {
	id := uuid.NewString()
	_, err := db.ExecContext(ctx, addEventQuery, raceID, id, description, eventType, distance)
	return id, err
}

func GetRaceEvents(ctx context.Context, db *db.DBLayer, raceID string) ([]EventObject, error) {
	var dbEvents []EventRow
	err := db.SelectContext(ctx, &dbEvents, listRaceEventsQuery, raceID)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	results := []EventObject{}
	for _, event := range dbEvents {
		results = append(results, EventObject{
			EventID:     event.EventID,
			Description: event.Description,
			Type:        event.Type,
			Distance:    event.Distance,
		})
	}

	return results, nil
}

func GetRaceEvent(ctx context.Context, db *db.DBLayer, raceID, eventID string) (EventObject, error) {
	var event EventRow

	err := db.GetContext(ctx, &event, selectEventQuery, raceID, eventID)
	if err != nil {
		return EventObject{}, err
	}

	return EventObject{
		EventID:     eventID,
		Description: event.Description,
		Type:        event.Type,
		Distance:    event.Distance,
	}, nil
}

func DeleteRaceEvent(ctx context.Context, db *db.DBLayer, raceID, eventID string) error {
	_, err := db.ExecContext(ctx, deleteEventQuery, raceID, eventID)
	if err != nil {
		return err
	}

	return nil
}

func GetRaceEventCount(ctx context.Context, db *db.DBLayer, raceID string) (int, error) {
	row := db.QueryRowContext(ctx, selectRaceEventCountQuery, raceID)
	if row.Err() != nil {
		return 0, row.Err()
	}

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
