package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/jmoiron/sqlx"
)

type EventResultsDao struct {
	db *sqlx.DB
}

type EventResultsDTO struct {
	FirstName sql.NullString `db:"first_name"`
	LastName  sql.NullString `db:"last_name"`
	BibNumber string         `db:"bib_number"`
	BirthYear sql.NullInt32  `db:"birth_year"`
	Gender    sql.NullString `db:"gender"`
	Result    int            `db:"result"`
}

func NewEventResultsDAO() *EventResultsDao {
	return &EventResultsDao{
		db: getDBInstance(),
	}
}

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

func (er *EventResultsDao) GetEventResults(ctx context.Context) ([]EventResultsDTO, error) {
	var results []EventResultsDTO
	err := getDBInstance().Select(&results, getEventResults, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx))
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (er *EventResultsDao) RecordFinisher(ctx context.Context, bib string) error {
	result, err := er.db.Exec(recordFinisherQuery, bib, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx))
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("no unclaimed times found")
	}

	return nil
}
