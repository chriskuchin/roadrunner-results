package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type RaceDAO struct {
	db     *sqlx.DB
	limit  int64
	offset int
}

type RaceDTO struct {
	RaceID  string `db:"race_id"`
	OwnerID string `db:"owner_id"`
	Name    string `db:"race_name"`
}

func NewRaceDAO(db *sqlx.DB) *RaceDAO {

	return &RaceDAO{
		db:    db,
		limit: 10,
	}
}

func (r *RaceDAO) CreateRace(ctx context.Context, id, name string) error {
	_, err := r.db.Exec("insert into races (race_id, race_name, owner_id) VALUES(?, ?, ?)", id, name, "123")
	return err
}

func (r *RaceDAO) GetRaceByID(ctx context.Context, id string) (RaceDTO, error) {
	dto := []RaceDTO{}
	err := r.db.Select(&dto, "select * from races where race_id = ?", id)
	if err != nil {
		return RaceDTO{}, err
	}
	log.Debug().Msgf("%+v", dto)

	return dto[0], nil
}

func (r *RaceDAO) ListRaces(ctx context.Context) ([]RaceDTO, error) {
	dto := []RaceDTO{}
	err := r.db.Select(&dto, "select * from races LIMIT ?", r.limit, r.offset)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}
	log.Debug().Msgf("%+v - %d", dto, len(dto))

	return dto, nil
}

func (r *RaceDAO) DeleteRace(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "delete from races where race_id = ?", id)
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}
	return nil
}
