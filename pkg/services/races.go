package services

import (
	"context"

	dao "github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/rs/zerolog/log"
)

var instance *RaceService

type RaceService struct {
	raceDAO dao.RaceDAO
}

type RaceResult struct {
	Name    string
	ID      string
	OwnerID string
}

func NewRaceService(raceDAO dao.RaceDAO) {
	if instance == nil {
		instance = &RaceService{
			raceDAO: raceDAO,
		}
	}
}

func GetRaceServiceInstance() *RaceService {
	return instance
}

func (rs *RaceService) CreateRace(ctx context.Context, name string) (string, error) {
	id, err := rs.raceDAO.CreateRace(ctx, name)
	if err != nil {
		log.Error().Err(err).Send()
		return "", err
	}

	return id, nil
}

func (rs *RaceService) ListRaces(ctx context.Context) ([]RaceResult, error) {
	races, err := rs.raceDAO.ListRaces(ctx)
	if err != nil {
		return nil, err
	}

	result := []RaceResult{}
	for _, race := range races {
		result = append(result, RaceResult{
			Name:    race.Name,
			ID:      race.RaceID,
			OwnerID: race.OwnerID,
		})
	}
	return result, nil
}

func (rs *RaceService) DeleteRace(ctx context.Context, id string) error {
	return rs.raceDAO.DeleteRace(ctx, id)
}
