package services

import (
	"context"

	dao "github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/google/uuid"
)

var raceInstance *RaceService

type RaceService struct {
	raceDAO *dao.RaceDAO
}

type RaceResult struct {
	Name    string
	ID      string
	OwnerID string
}

func NewRaceService(raceDAO *dao.RaceDAO) {
	if raceInstance == nil {
		raceInstance = &RaceService{
			raceDAO: raceDAO,
		}
	}
}

func GetRaceServiceInstance() *RaceService {
	return raceInstance
}

func (rs *RaceService) CreateRace(ctx context.Context, name string) (string, error) {
	id := uuid.NewString()

	return id, rs.CreateRaceWithID(ctx, id, name)
}

func (rs *RaceService) CreateRaceWithID(ctx context.Context, id, name string) error {
	return rs.raceDAO.CreateRace(ctx, id, name)
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
