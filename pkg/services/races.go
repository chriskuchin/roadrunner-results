package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

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

func GetRace(ctx context.Context, db *sqlx.DB) (RaceResult, error) {
	dto := []RaceDTO{}
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
	_, err := db.Exec(createRaceQuery, id, name, "123")
	return err
}

func ListRaces(ctx context.Context, db *sqlx.DB) ([]RaceResult, error) {
	races := []RaceDTO{}
	err := db.Select(&races, "select * from races LIMIT ?", 20)
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

func DeleteRace(ctx context.Context, db *sqlx.DB, id string) error {
	_, err := db.ExecContext(ctx, "delete from races where race_id = ?", id)
	return err
}

const (
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
	selectRaceParticipantCountQuery = `
		select count(1) from participants
			where
				race_id = ?
	`
	selectRaceParticipantGenderCountQuery = `
		select gender, count(1) from participants
			where race_id = ?
			GROUP BY gender
	`
	selectRaceParticipantBirthYearCountQuery = `
		select birth_year, gender, count(1) from participants
			where race_id = ?
			GROUP BY birth_year, gender
			ORDER BY birth_year
	`
)

type RaceDTO struct {
	Name    string `db:"race_name"`
	ID      string `db:"race_id"`
	OwnerID string `db:"owner_id"`
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

func GetRaceParticipantCount(ctx context.Context, db *sqlx.DB) (int, error) {
	row := db.QueryRow(selectRaceParticipantCountQuery, util.GetRaceIDFromContext(ctx))
	if row.Err() != nil {
		return 0, row.Err()
	}

	var count int
	err := row.Scan(&count)
	return count, err
}

func GetRaceGenderCount(ctx context.Context, db *sqlx.DB) (female int, male int, err error) {
	rows, err := db.Query(selectRaceParticipantGenderCountQuery, util.GetRaceIDFromContext(ctx))
	if err != nil {
		return
	}

	for rows.Next() {
		var gender string
		var count int
		rows.Scan(&gender, &count)
		if gender == "F" {
			female = count
		} else if gender == "M" {
			male = count
		}
	}

	return
}

func GetRaceParticipantsBirthYearCount(ctx context.Context, db *sqlx.DB) ([]map[string]interface{}, error) {
	rows, err := db.Query(selectRaceParticipantBirthYearCountQuery, util.GetRaceIDFromContext(ctx))
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{} = []map[string]interface{}{}
	for rows.Next() {
		var year int
		var count int
		var gender string
		rows.Scan(&year, &gender, &count)
		results = append(results, map[string]interface{}{
			"year":   year,
			"gender": gender,
			"count":  count,
		})
	}

	return results, nil
}
