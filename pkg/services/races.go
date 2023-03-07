package services

import (
	"context"

	dao "github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var raceInstance *RaceService

type RaceService struct {
	raceDAO *dao.RaceDAO
}

type RaceResult struct {
	Name             string `json:"name"`
	ID               string `json:"id"`
	OwnerID          string `json:"owner_id"`
	EventCount       int    `json:"event_count"`
	ParticipantCount int    `json:"participant_count"`
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

func (rs *RaceService) GetRace(ctx context.Context) (RaceResult, error) {
	race, err := rs.raceDAO.GetRaceByID(ctx, util.GetRaceIDFromContext(ctx))
	if err != nil {
		return RaceResult{}, err
	}

	eventCount, err := rs.raceDAO.GetRaceEventCount(ctx)
	if err != nil {
		log.Error().Err(err).Send()
		return RaceResult{}, err
	}

	participantCount, err := rs.raceDAO.GetRaceParticipantCount(ctx)
	if err != nil {
		log.Error().Err(err).Send()
		return RaceResult{}, err
	}

	return RaceResult{
		Name:             race.Name,
		ID:               race.ID,
		OwnerID:          race.OwnerID,
		EventCount:       eventCount,
		ParticipantCount: participantCount,
	}, nil
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
			ID:      race.ID,
			OwnerID: race.OwnerID,
		})
	}
	return result, nil
}

func (rs *RaceService) DeleteRace(ctx context.Context, id string) error {
	return rs.raceDAO.DeleteRace(ctx, id)
}
