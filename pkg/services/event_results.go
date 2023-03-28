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
		BibNumber string         `db:"bib_number"`
		BirthYear sql.NullInt32  `db:"birth_year"`
		Gender    sql.NullString `db:"gender"`
		Result    int            `db:"result"`
	}

	ParticipantEventResult struct {
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
		BibNumber string `json:"bib_number"`
		BirthYear int    `json:"birth_year,omitempty"`
		Gender    string `json:"gender,omitempty"`
		Result    int    `json:"result_ms"`
	}
)

const (
	getEventResults string = `
		SELECT
			p.first_name,
			p.last_name,
			r.bib_number,
			p.birth_year,
			p.gender,
			r.result
		FROM results as r
		LEFT JOIN participants as p
			USING(bib_number, race_id)
		WHERE
			race_id = ?
			AND
			event_id = ?
		ORDER by p.birth_year, p.gender, r.result;
	`

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

func GetEventResults(ctx context.Context, db *sqlx.DB) ([]ParticipantEventResult, error) {
	var results []EventResultsRow
	err := db.Select(&results, getEventResults, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx))
	if err != nil {
		return nil, err
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	var participantResults []ParticipantEventResult = []ParticipantEventResult{}
	for _, result := range results {
		participantResults = append(participantResults, ParticipantEventResult{
			BibNumber: result.BibNumber,
			Result:    result.Result,
			FirstName: result.FirstName.String,
			LastName:  result.LastName.String,
			BirthYear: int(result.BirthYear.Int32),
			Gender:    result.Gender.String,
		})
	}

	return participantResults, nil
}

func RecordTimerResult(ctx context.Context, endTS int64) error {
	start, _, err := GetActiveTimerStart(ctx)
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}

	return InsertPartialResult(ctx, endTS-start)
}

func RecordFinisherResult(ctx context.Context, db *sqlx.DB, bib string) error {
	result, err := db.Exec(recordFinisherQuery, bib, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx))
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("no unclaimed times found")
	}

	return err
}
