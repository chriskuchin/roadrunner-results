package services

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type AttemptsRow struct {
	RaceID  string  `db:"race_id"`
	EventID string  `db:"event_id"`
	Bib     string  `db:"bib"`
	Attempt int     `db:"attempt_no"`
	Result  float32 `db:"result"`
}

const (
	insertAttemptQuery = `INSERT into attempts (race_id, event_id, bib, attempt_no, result) VALUES (?, ?, ?, ?, ?)`

	upsertResultQuery = `INSERT into results(race_id, event_id, bib_number, result)
    SELECT race_id, event_id, bib as bib_number, MAX(result) from attempts where race_id = ? AND event_id = ? AND bib = ?
    ON CONFLICT(race_id, event_id, bib_number) DO UPDATE SET result=excluded.result`
	//  INSERT INTO results(race_id, event_id, bib_number, result) SELECT race_id, event_id, bib as bib_number, MAX(result) from attempts
	listAttemptsQuery = `SELECT * FROM attempts where race_id = ? AND event_id = ? AND bib = ? ORDER BY result DESC`
)

func RecordAttempt(ctx context.Context, db *sqlx.DB, raceID, eventID, bib string, attemptNo int, result float32) error {
	_, err := db.Exec(insertAttemptQuery, raceID, eventID, bib, attemptNo, result)
	return err
}

func ListAttempts(ctx context.Context, db *sqlx.DB, raceID, eventID, bib string) error {
	var attempts []AttemptsRow
	err := db.Select(&attempts, listAttemptsQuery, raceID, eventID, bib)
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}

	for _, attempt := range attempts {
		fmt.Println(attempt)
	}
	return nil
}

func UpsertBestAttemptResults(ctx context.Context, db *sqlx.DB, raceID, eventID, bib string) error {
	rslt, err := db.Exec(upsertResultQuery, raceID, eventID, bib)
	if err != nil {
		log.Err(err).Send()
		return err
	}

	rows, _ := rslt.RowsAffected()
	log.Info().Int64("result", rows).Send()
	return nil
}
