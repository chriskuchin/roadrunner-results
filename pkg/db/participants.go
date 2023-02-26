package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type (
	ParticipantsDAO struct {
		db *sqlx.DB
	}

	Participant struct {
		RaceID    string `db:"race_id"`
		EventID   string `db:"event_id"`
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
			event_id,
			bib_number,
			first_name,
			last_name,
			gender,
			team,
			birth_year
		)
		VALUES (
			:race_id,
			:event_id,
			:bib_number,
			:first_name,
			:last_name,
			:gender,
			:team,
			:birth_year
		)
	`
)

func NewParticipantsDAO(db *sqlx.DB) *ParticipantsDAO {
	return &ParticipantsDAO{
		db: db,
	}
}

func (dao *ParticipantsDAO) InsertParticipant(ctx context.Context, participant Participant) error {
	_, err := dao.db.NamedExec(addParticipantQuery, participant)
	return err
}
