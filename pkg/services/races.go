package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var raceInstance *RaceService

type RaceService struct {
	raceDAO *db.RaceDAO
}

type RaceResult struct {
	Name             string           `json:"name"`
	ID               string           `json:"id"`
	OwnerID          string           `json:"owner_id"`
	EventCount       int              `json:"event_count"`
	Events           []EventObject    `json:"events"`
	ParticipantStats ParticipantStats `json:"participant_stats"`
}

type ParticipantStats struct {
	MaleCount          int                      `json:"male"`
	FemaleCount        int                      `json:"female"`
	Total              int                      `json:"total"`
	BirthYearHistogram []map[string]interface{} `json:"birth_year_distribution"`
}

func NewRaceService() {
	if raceInstance == nil {
		raceInstance = &RaceService{
			raceDAO: db.NewRaceDAO(),
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
	}

	participantCount, err := rs.raceDAO.GetRaceParticipantCount(ctx)
	if err != nil {
		log.Error().Err(err).Send()
	}

	femaleCount, maleCount, err := rs.raceDAO.GetRaceGenderCount(ctx)
	if err != nil {
		log.Error().Err(err).Send()
	}

	birthYearDistro, err := rs.raceDAO.GetRaceParticipantsBirthYearCount(ctx)
	if err != nil {
		log.Error().Err(err).Send()
	}

	events, err := db.NewEventDAO().GetRaceEvents(ctx)
	if err != nil {
		log.Error().Err(err).Send()
	}

	eventResults := []EventObject{}
	for _, event := range events {
		eventResults = append(eventResults, EventObject{
			Description: event.Description,
			EventID:     event.EventID,
			Distance:    event.Distance,
		})
	}

	return RaceResult{
		Name:       race.Name,
		ID:         race.ID,
		OwnerID:    race.OwnerID,
		EventCount: eventCount,
		Events:     eventResults,
		ParticipantStats: ParticipantStats{
			Total:              participantCount,
			FemaleCount:        femaleCount,
			MaleCount:          maleCount,
			BirthYearHistogram: birthYearDistro,
		},
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
