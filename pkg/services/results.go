package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/jmoiron/sqlx"
)

const (
	addResultQuery = `
		insert into results (
			race_id,
			event_id,
			bib_number,
			result
		)
		VALUES (?, ?, ?, ?)
	`

	insertPartialQuery string = `
		insert into results (
			race_id,
			event_id,
			bib_number,
			result
		)
		VALUES(?, ?, NULL, ?)
	`
)

func InsertResults(ctx context.Context, db *sqlx.DB, raceID, eventID, bibNumber, result string) error {
	_, err := db.Exec(addResultQuery, raceID, eventID, bibNumber, result)
	return err
}

// select p.first_name, p.last_name, p.bib_number, r.result from participants as p join results as r on p.race_id = r.race_id and p.event_id = r.event_id and p.bib_number = r.bib_number
// select p.first_name, p.last_name, p.bib_number, p.birth_year, r.result from participants as p join results as r on p.race_id = r.race_id and p.event_id = r.event_id and p.bib_number = r.bib_number group by p.birth_year order by r.result;

func InsertPartialResult(ctx context.Context, db *sqlx.DB, result int64) error {
	_, err := db.Exec(insertPartialQuery, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), result)
	return err
}
