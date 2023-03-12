package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type (
	ParticipantsDAO struct {
		db *sqlx.DB
	}

	Participant struct {
		RaceID    string `db:"race_id"`
		BibNumber string `db:"bib_number"`
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		BirthYear int    `db:"birth_year"`
		Gender    string `db:"gender"`
		Team      string `db:"team"`
	}
)

const (
	addParticipantQuery = `
		insert into participants (
			race_id,
			bib_number,
			first_name,
			last_name,
			gender,
			team,
			birth_year
		)
		VALUES (
			:race_id,
			:bib_number,
			:first_name,
			:last_name,
			:gender,
			:team,
			:birth_year
		)
	`

	listParticipantsQuery string = `
		select * from participants
			where
				race_id = ?
			limit ? offset ?
	`
)

func NewParticipantsDAO() *ParticipantsDAO {
	return &ParticipantsDAO{
		db: getDBInstance(),
	}
}

func (dao *ParticipantsDAO) InsertParticipant(ctx context.Context, participant Participant) error {
	_, err := dao.db.NamedExec(addParticipantQuery, participant)
	return err
}

func (dao *ParticipantsDAO) ListParticipants(ctx context.Context, raceID string, limit, offset int) ([]Participant, error) {
	log.Info().Int("limit", limit).Int("offset", offset).Send()
	var results []Participant
	return results, dao.db.Select(&results, listParticipantsQuery, raceID, limit, offset)
}
