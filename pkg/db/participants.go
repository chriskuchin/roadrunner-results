package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type (
	ParticipantsDAO struct {
		db sqlx.DB
	}

	Participant struct {
		RaceID    string `db:"race_id"`
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		Team      string `db:"team"`
		BirthYear int    `db:"birth_year"`
	}
)

func NewParticipantsDAO(db sqlx.DB) *ParticipantsDAO {
	return &ParticipantsDAO{
		db: db,
	}
}

func (dao *ParticipantsDAO) BatchInsertParticipants(ctx context.Context, participants []Participant) error {
	_, err := dao.db.NamedExec("insert into participants (race_id, first_name, last_name, team, birth_year) VALUES (:race_id, :first_name, :last_name, :team, :birth_year)", participants)
	log.Info().Err(err).Send()
	return err
}
