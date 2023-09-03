package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const (
	getRaceOwnerQuery string = `
		select owner_id from races
			where race_id = ?
	`

	createRaceQuery string = `
		insert into races (
			race_id,
			race_name,
			owner_id
		)
		VALUES(?, ?, ?)
	`
	selectRaceEventCountQuery string = `
		select count(1) from events
			where
				race_id = ?
	`
)

type RaceRow struct {
	Name    string `db:"race_name"`
	ID      string `db:"race_id"`
	OwnerID string `db:"owner_id"`
}

type RaceResult struct {
	Name             string            `json:"name"`
	ID               string            `json:"id"`
	OwnerID          string            `json:"owner_id"`
	EventCount       int               `json:"event_count,omitempty"`
	Events           []EventObject     `json:"events,omitempty"`
	ParticipantStats *ParticipantStats `json:"participant_stats,omitempty"`
}

type ParticipantStats struct {
	MaleCount          int                      `json:"male,omitempty"`
	FemaleCount        int                      `json:"female,omitempty"`
	Total              int                      `json:"total,omitempty"`
	BirthYearHistogram []map[string]interface{} `json:"birth_year_distribution,omitempty"`
}

func GetRaceOwnerID(ctx context.Context, db *sqlx.DB, id string) (string, error) {
	var ownerID []string
	err := db.Select(&ownerID, getRaceOwnerQuery, id)
	if err != nil {
		return "", err
	}

	if len(ownerID) != 1 {
		return "", fmt.Errorf("unexpected number of ownerIDs returned")
	}

	return ownerID[0], nil
}

func GetRace(ctx context.Context, db *sqlx.DB) (RaceResult, error) {
	dto := []RaceRow{}
	err := db.Select(&dto, "select * from races where race_id = ?", util.GetRaceIDFromContext(ctx))
	if err != nil {
		return RaceResult{}, err
	}

	race := dto[0]
	if err != nil {
		return RaceResult{}, err
	}

	eventCount, err := GetRaceEventCount(ctx, db)
	if err != nil {
		log.Error().Err(err).Send()
	}

	participantCount, err := GetRaceParticipantCount(ctx, db)
	if err != nil {
		log.Error().Err(err).Send()
	}

	femaleCount, maleCount, err := GetRaceGenderCount(ctx, db)
	if err != nil {
		log.Error().Err(err).Send()
	}

	birthYearDistro, err := GetRaceParticipantsBirthYearCount(ctx, db)
	if err != nil {
		log.Error().Err(err).Send()
	}

	events, err := GetRaceEvents(ctx, db)
	if err != nil {
		log.Error().Err(err).Send()
	}

	eventResults := []EventObject{}
	for _, event := range events {
		eventResults = append(eventResults, EventObject{
			Description: event.Description,
			EventID:     event.EventID,
			Type:        event.Type,
			Distance:    event.Distance,
		})
	}

	return RaceResult{
		Name:       race.Name,
		ID:         race.ID,
		OwnerID:    race.OwnerID,
		EventCount: eventCount,
		Events:     eventResults,
		ParticipantStats: &ParticipantStats{
			Total:              participantCount,
			FemaleCount:        femaleCount,
			MaleCount:          maleCount,
			BirthYearHistogram: birthYearDistro,
		},
	}, nil
}

func CreateRace(ctx context.Context, db *sqlx.DB, name string) (string, error) {
	id := uuid.NewString()
	err := CreateRaceWithID(ctx, db, id, name)
	return id, err
}

func CreateRaceWithID(ctx context.Context, db *sqlx.DB, id, name string) error {
	_, err := db.Exec(createRaceQuery, id, name, util.GetCurrentUserID(ctx))
	return err
}

func ListRaces(ctx context.Context, db *sqlx.DB) ([]RaceResult, error) {
	races := []RaceRow{}
	err := db.Select(&races, "select * from races ORDER BY rowid DESC LIMIT ?", 20)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	result := []RaceResult{}
	for _, race := range races {
		result = append(result, RaceResult{
			Name:    race.Name,
			ID:      race.ID,
			OwnerID: race.OwnerID,
		})
	}
	return result, nil
}

var raceTables []string = []string{
	"races",
	"events",
	"participants",
	"event_participation",
	"results",
	"heats",
	"timers",
	"timer_results",
}

func DeleteRace(ctx context.Context, db *sqlx.DB, id string) (err error) {
	var failedObjects []string = []string{}
	for _, table := range raceTables {
		_, err := db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE race_id = ?", table), id)
		if err != nil {
			log.Error().Err(err).Send()
			failedObjects = append(failedObjects, table)
		}
	}

	if len(failedObjects) > 0 {
		err = fmt.Errorf("failed deleting race objects: %s", strings.Join(failedObjects, ", "))
	}

	return err
}
