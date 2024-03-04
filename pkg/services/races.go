package services

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chriskuchin/roadrunner-results/pkg/client"
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
			owner_id,
			race_date
		)
		VALUES(?, ?, ?, ?)
	`
	selectRaceEventCountQuery string = `
		select count(1) from events
			where
				race_id = ?
	`
)

type RaceRow struct {
	Name     string        `db:"race_name"`
	ID       string        `db:"race_id"`
	OwnerID  string        `db:"owner_id"`
	RaceDate sql.NullInt64 `db:"race_date"`
}

type RaceResult struct {
	Name             string            `json:"name"`
	Date             *string           `json:"date,omitempty"`
	ID               string            `json:"id"`
	OwnerID          string            `json:"owner_id"`
	EventCount       int               `json:"event_count,omitempty"`
	Events           []EventStats      `json:"events,omitempty"`
	ParticipantStats *ParticipantStats `json:"participant_stats,omitempty"`
}

type EventStats struct {
	EventID     string `json:"eventId"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Distance    int    `json:"distance"`
	Finishers   int    `json:"finishers"`
}

type ParticipantStats struct {
	MaleCount          int                      `json:"male,omitempty"`
	FemaleCount        int                      `json:"female,omitempty"`
	Total              int                      `json:"total,omitempty"`
	Finishers          int                      `json:"finishers,omitempty"`
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

func GetRace(ctx context.Context, db *sqlx.DB, raceID string) (RaceResult, error) {
	dto := []RaceRow{}
	err := db.Select(&dto, "select * from races where race_id = ?", raceID)
	if err != nil {
		return RaceResult{}, err
	}

	race := dto[0]

	eventCount, err := GetRaceEventCount(ctx, db, raceID)
	if err != nil {
		log.Error().Err(err).Send()
	}

	participantCount, err := GetRaceParticipantCount(ctx, db, raceID)
	if err != nil {
		log.Error().Err(err).Send()
	}

	finisherCount, err := GetRaceFinisherCount(ctx, db, raceID)
	if err != nil {
		log.Error().Err(err).Send()
	}

	femaleCount, maleCount, err := GetRaceGenderCount(ctx, db, raceID)
	if err != nil {
		log.Error().Err(err).Send()
	}

	birthYearDistro, err := GetRaceParticipantsBirthYearCount(ctx, db, raceID)
	if err != nil {
		log.Error().Err(err).Send()
	}

	events, err := GetRaceEvents(ctx, db, raceID)
	if err != nil {
		log.Error().Err(err).Send()
	}

	eventResults := []EventStats{}
	for _, event := range events {
		finisherCount, err := GetEventFinisherCount(ctx, db, raceID, event.EventID)
		if err != nil {
			log.Error().Err(err).Send()
		}
		eventResults = append(eventResults, EventStats{
			Description: event.Description,
			EventID:     event.EventID,
			Type:        event.Type,
			Distance:    event.Distance,
			Finishers:   finisherCount,
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
			Finishers:          finisherCount,
			FemaleCount:        femaleCount,
			MaleCount:          maleCount,
			BirthYearHistogram: birthYearDistro,
		},
	}, nil
}

func ImportRaceFromURL(ctx context.Context, db *sqlx.DB, url string, date time.Time, name string) (string, error) {
	eventInfo := client.GetEventInformation(ctx, url)
	raceID, err := CreateRace(ctx, db, name, date)
	ctx = context.WithValue(ctx, util.RaceID, raceID)
	if err != nil {
		return "", err
	}

	for _, event := range eventInfo {
		eventID, err := AddEvent(ctx, db, raceID, event.Name, "", int(event.Distance))
		if err != nil {
			return "", err
		}

		ctx = context.WithValue(ctx, util.EventID, eventID)
		for _, result := range event.Results {
			pRow := ParticipantRow{
				RaceID:    raceID,
				BibNumber: result.Bib,
				FirstName: result.FirstName,
				LastName:  result.LastName,
				Gender:    result.Gender,
				Team:      result.Team,
			}

			if result.Grade != "" {
				grade, _ := strconv.Atoi(result.Grade)
				pRow.Grade = sql.NullInt64{Int64: int64(grade), Valid: true}
			}

			AddParticipant(ctx, db, pRow)

			err = InsertResults(ctx, db, raceID, eventID, result.Bib, fmt.Sprintf("%d", result.Time))
			if err != nil {
				log.Error().Err(err).Send()
			}
		}
	}

	return raceID, nil
}

func CreateRace(ctx context.Context, db *sqlx.DB, name string, date time.Time) (string, error) {
	id := uuid.NewString()
	err := CreateRaceWithID(ctx, db, id, name, date)
	return id, err
}

func CreateRaceWithID(ctx context.Context, db *sqlx.DB, id, name string, date time.Time) error {
	dateNum, err := strconv.Atoi(date.Format("20060102"))
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}
	_, err = db.Exec(createRaceQuery, id, name, util.GetCurrentUserID(ctx), dateNum)
	return err
}

func ListRaces(ctx context.Context, db *sqlx.DB, limit, offset int) ([]RaceResult, error) {
	races := []RaceRow{}
	err := db.Select(&races, "select * from races ORDER BY race_date DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	result := []RaceResult{}
	for _, race := range races {
		entry := RaceResult{
			Name:    race.Name,
			ID:      race.ID,
			OwnerID: race.OwnerID,
		}
		if race.RaceDate.Valid {
			date, err := time.Parse("20060102", fmt.Sprintf("%d", race.RaceDate.Int64))
			if err == nil {
				formatted := date.Format("2006-01-02")
				entry.Date = &formatted
			}
		}
		result = append(result, entry)
	}
	return result, nil
}

var raceTables []string = []string{
	"races",
	"events",
	"participants",
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
