package db

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type RaceDAO struct {
	db *sqlx.DB
}

type RaceDTO struct {
	RaceID  string `db:"race_id"`
	OwnerID string `db:"owner_id"`
	Race    string `db:"race"`
}

func NewRaceDAO(db *sqlx.DB) *RaceDAO {

	return &RaceDAO{
		db: db,
	}
}

func (r *RaceDAO) CreateRace(name string) (string, error) {
	id := uuid.New()
	_, err := r.db.Exec("insert into races (race_id, race, owner_id) VALUES(?, ?, ?)", id, name, "123")
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (r *RaceDAO) GetRaceByID(id string) {
	dto := []RaceDTO{}
	r.db.Select(&dto, "select * from races where race_id = ?", id)
	log.Debug().Msgf("%+v", dto)
}
