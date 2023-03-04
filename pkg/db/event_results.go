package db

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type EventResultsDao struct {
	db *sqlx.DB
}

type EventResultsDTO struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	BibNumber string `db:"bib_number"`
	BirthYear int    `db:"birth_year"`
	Gender    string `db:"gender"`
	Result    int    `db:"result"`
}

func NewEventResultsDAO(db *sqlx.DB) *EventResultsDao {
	return &EventResultsDao{
		db: db,
	}
}

const (
	getEventResults string = `
		SELECT
			p.first_name,
			p.last_name,
			p.bib_number,
			p.birth_year,
			p.gender,
			r.result
		FROM participants as p
		JOIN results as r
			USING(bib_number, race_id)
		WHERE
			p.bib_number in (
				SELECT bib_number
				FROM event_participation
				WHERE
					race_id = ?
					AND
					event_id = ?
			)
		ORDER by p.birth_year, p.gender, r.result;
	`
)

func (er *EventResultsDao) GetEventResults(ctx context.Context) error {
	var results []EventResultsDTO
	err := er.db.Select(&results, getEventResults, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx))
	if err != nil {
		return err
	}

	for _, result := range results {
		log.Info().Msgf("%+v", result)
	}

	return nil
}
