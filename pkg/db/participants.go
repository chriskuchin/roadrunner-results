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
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		BirthYear int    `db:"birth_year"`
		Gender    string `db:"gender"`
	}
)

func NewParticipantsDAO(db *sqlx.DB) *ParticipantsDAO {
	return &ParticipantsDAO{
		db: db,
	}
}

func (dao *ParticipantsDAO) BatchInsertParticipants(ctx context.Context, participants []Participant) error {
	_, err := dao.db.NamedExec("insert into participants (race_id, first_name, last_name, team, birth_year) VALUES (:race_id, :first_name, :last_name, :team, :birth_year)", participants)
	return err
}
