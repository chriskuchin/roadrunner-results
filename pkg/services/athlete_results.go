package services

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type AthleteResultsRow struct {
	RaceID    string `db:"race_id"`
	EventID   string `db:"event_id"`
	TimerID   string `db:"timer_id"`
	BibNumber string `db:"bib_number"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Gender    string `db:"gender"`
	Team      string `db:"team"`
	BirthYear string `db:"birth_year"`
	Result    int    `db:"result"`
}

type AthleteResult struct {
	RaceID    string      `json:"race_id"`
	Event     EventObject `json:"event"`
	BibNumber string      `json:"bib_number"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Gender    string      `json:"gender"`
	Team      string      `json:"team"`
	BirthYear string      `json:"birth_year"`
	Result    int         `json:"result"`
}

func FindAthleteResults(ctx context.Context, db *sqlx.DB, filters map[string]string) ([]AthleteResult, error) {
	if len(filters) == 0 {
		return nil, fmt.Errorf("invalid filters")
	}

	var query string = "SELECT * FROM results JOIN participants USING(race_id, bib_number) WHERE"
	var first bool = true
	var values []interface{} = []interface{}{}
	for k, v := range filters {
		if !first {
			query = fmt.Sprintf("%s AND", query)
		}

		if k == "gender" {
			query = fmt.Sprintf("%s gender = ?", query)
			values = append(values, v)
		} else if k == "first_name" {
			query = fmt.Sprintf("%s first_name LIKE ?", query)
			values = append(values, v)
		} else if k == "last_name" {
			query = fmt.Sprintf("%s last_name LIKE ?", query)
			values = append(values, v)
		} else if k == "birth_year" {
			query = fmt.Sprintf("%s birth_year = ?", query)
			values = append(values, v)
		}

		first = false
	}

	var results []AthleteResultsRow
	err := db.Select(&results, query, values...)
	if err != nil {
		return nil, err
	}

	var athleteResults []AthleteResult = []AthleteResult{}
	for _, row := range results {

		event, err := GetRaceEvent(ctx, db, row.RaceID, row.EventID)
		if err != nil {
			log.Error().Err(err).Send()
			continue
		}
		log.Info().Interface("result", row).Interface("event", event).Err(err).Send()

		athleteResults = append(athleteResults, AthleteResult{
			Event:     event,
			RaceID:    row.RaceID,
			BibNumber: row.BibNumber,
			FirstName: row.FirstName,
			LastName:  row.LastName,
			Gender:    row.Gender,
			Team:      row.Team,
			BirthYear: row.BirthYear,
			Result:    row.Result,
		})
	}

	return athleteResults, nil
}
