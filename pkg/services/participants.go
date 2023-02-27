package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
)

var participantServiceInstance *ParticipantsService

type (
	ParticipantsService struct {
		participantsDAO *db.ParticipantsDAO
	}

	ParticipantObject struct {
		BibNumber string `json:"bibNumber"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		BirthYear int    `json:"birthYear"`
		Gender    string `json:"gender"`
		Team      string `json:"team"`
	}
)

func NewParticipantsService(db *db.ParticipantsDAO) {
	if participantServiceInstance == nil {
		participantServiceInstance = &ParticipantsService{
			participantsDAO: db,
		}
	}
}

func GetParticipantServiceInstance() *ParticipantsService {
	return participantServiceInstance
}

func (ps *ParticipantsService) AddParticipant(ctx context.Context, participant db.Participant) error {
	return ps.participantsDAO.InsertParticipant(ctx, participant)
}

func (ps *ParticipantsService) ListParticipants(ctx context.Context, limit, offset int) ([]ParticipantObject, error) {
	dbResults, err := ps.participantsDAO.ListParticipants(ctx, util.GetRaceIDFromContext(ctx), limit, offset)
	if err != nil {
		return nil, err
	}

	results := []ParticipantObject{}
	for _, participant := range dbResults {
		results = append(results, ParticipantObject{
			BibNumber: participant.BibNumber,
			FirstName: participant.FirstName,
			LastName:  participant.LastName,
			BirthYear: participant.BirthYear,
			Gender:    participant.Gender,
			Team:      participant.Team,
		})
	}

	return results, nil
}
