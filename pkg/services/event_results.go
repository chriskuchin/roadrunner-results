package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type (
	EventResultsRow struct {
		FirstName sql.NullString `db:"first_name"`
		LastName  sql.NullString `db:"last_name"`
		BibNumber sql.NullString `db:"bib_number"`
		BirthYear sql.NullInt32  `db:"birth_year"`
		Gender    sql.NullString `db:"gender"`
		Result    int            `db:"result"`
		TimerID   sql.NullString `db:"timer_id"`
		Team      sql.NullString `db:"team"`
	}

	ParticipantEventResult struct {
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
		BibNumber string `json:"bib_number"`
		BirthYear int    `json:"birth_year,omitempty"`
		Gender    string `json:"gender,omitempty"`
		Result    int    `json:"result_ms"`
		TimerID   string `json:"timer_id"`
		Team      string `json:"team"`
	}
)

const (
	getEventResultsQuery string = "SELECT p.first_name, p.last_name, r.bib_number, p.birth_year, p.gender, r.result, r.timer_id, p.team FROM results as r LEFT JOIN participants as p USING(bib_number, race_id) WHERE"

	recordFinisherQuery string = `
		UPDATE results SET bib_number = ?
		WHERE rowid in (
			SELECT rowid from results
			WHERE
				bib_number is NULL
				AND
				race_id = ?
				AND
				event_id = ?
			ORDER by rowid LIMIT 1
		)
	`
)

func GetEventResults(ctx context.Context, db *sqlx.DB, filters map[string][]string) ([]ParticipantEventResult, error) {
	query := getEventResultsQuery

	var filterValues []interface{} = []interface{}{}
	isFirst := true
	for filter, values := range filters {
		if !isFirst {
			query = fmt.Sprintf("%s AND", query)
		} else {
			isFirst = false
		}
		query = fmt.Sprintf("%s (", query)
		for idx, val := range values {
			if idx != 0 {
				query = fmt.Sprintf("%s OR", query)
			}

			if filter == "timer" {
				query = fmt.Sprintf("%s r.timer_id = ?", query)
			} else if filter == "race" {
				query = fmt.Sprintf("%s r.race_id = ?", query)
			} else if filter == "gender" {
				query = fmt.Sprintf("%s p.gender = ?", query)
			} else if filter == "event" {
				query = fmt.Sprintf("%s r.event_id = ?", query)
			} else if filter == "team" {
				query = fmt.Sprintf("%s p.team = ?", query)
			} else if filter == "year" {
				query = fmt.Sprintf("%s p.birth_year = ?", query)
			} else if filter == "name" {
				query = fmt.Sprintf("%s first_name LIKE ? OR last_name LIKE ?", query)
				filterValues = append(filterValues, fmt.Sprintf("%%%s", val))
				filterValues = append(filterValues, fmt.Sprintf("%%%s", val))
			}

			if filter != "name" {
				filterValues = append(filterValues, val)
			}
		}
		query = fmt.Sprintf("%s )", query)
	}

	log.Debug().Str("query", query).Interface("filter", filterValues).Send()
	return runEventResultsQuery(db, query, filterValues...)
}

func runEventResultsQuery(db *sqlx.DB, query string, args ...interface{}) ([]ParticipantEventResult, error) {
	var results []EventResultsRow
	err := db.Select(&results, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Query Failed")
		return nil, err
	}

	var participantResults []ParticipantEventResult = []ParticipantEventResult{}
	for _, result := range results {
		participantResults = append(participantResults, ParticipantEventResult{
			BibNumber: result.BibNumber.String,
			Result:    result.Result,
			FirstName: result.FirstName.String,
			LastName:  result.LastName.String,
			BirthYear: int(result.BirthYear.Int32),
			Gender:    result.Gender.String,
			TimerID:   result.TimerID.String,
			Team:      result.Team.String,
		})
	}

	return participantResults, nil
}

func RecordElapsedTimeResult(ctx context.Context, elapsed int64) error {
	return InsertPartialResult(ctx, elapsed, "")
}

func RecordTimerResult(ctx context.Context, endTS int64) error {
	start, timerID, err := GetActiveTimerStart(ctx)
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}

	return InsertPartialResult(ctx, endTS-start, timerID)
}

func RecordFinisherResult(ctx context.Context, db *sqlx.DB, bib string) error {
	result, err := db.Exec(recordFinisherQuery, bib, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx))
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("no unclaimed times found")
	}

	return err
}
