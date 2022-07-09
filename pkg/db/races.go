package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type RaceDAO struct {
	db *sqlx.DB
}

type RaceDTO struct {
	ID   int    `db:"rowid"`
	Race string `db:"race"`
}

func NewRaceDAO(db *sqlx.DB) *RaceDAO {

	return &RaceDAO{
		db: db,
	}
}

func (r *RaceDAO) CreateRace(name string) {
	_, err := r.db.Exec("insert into races (race) VALUES(:race)", name)

	log.Error().Err(err).Send()

	dto := []RaceDTO{}
	r.db.Select(&dto, "select * from races")
	log.Info().Msgf("%+v", dto)
}

func (r *RaceDAO) GetRaceByID(id int) {

}
