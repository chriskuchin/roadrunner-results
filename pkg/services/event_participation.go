package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
)

var eventParticipationInstance *EventParticipationService

type (
	EventParticipationService struct {
		eventParticipationDao *db.EventParticipationDao
	}
)

func NewEventParticipationService(db *db.EventParticipationDao) {
	if eventParticipationInstance == nil {
		eventParticipationInstance = &EventParticipationService{
			eventParticipationDao: db,
		}
	}
}

func GetEventParticipationInstance() *EventParticipationService {
	return eventParticipationInstance
}

func (ep *EventParticipationService) RecordEventParticipation(ctx context.Context, raceID, eventID, bibNumber string) error {
	return ep.eventParticipationDao.RecordEventParticipation(ctx, raceID, eventID, bibNumber)
}
