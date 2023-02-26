package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
)

var participantServiceInstance *ParticipantsService

type (
	ParticipantsService struct {
		participantsDAO *db.ParticipantsDAO
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
