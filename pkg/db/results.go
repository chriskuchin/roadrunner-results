package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ResultsDAO struct {
	db *sqlx.DB
}

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
)

func NewResultsDAO() *ResultsDAO {
	return &ResultsDAO{
		db: getDBInstance(),
	}
}

func (dao *ResultsDAO) InsertResult(ctx context.Context, raceID, eventID, bibNumber string, result int64) error {
	_, err := dao.db.Exec(addResultQuery, raceID, eventID, bibNumber, result)
	return err
}
