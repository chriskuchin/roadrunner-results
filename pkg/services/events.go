package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
)

var eventsInstance *EventsService

type EventsService struct {
	eventsDao db.EventDao
}

func NewEventService(dao db.EventDao) {
	if eventsInstance == nil {
		eventsInstance = &EventsService{
			eventsDao: dao,
		}
	}
}

func GetEventsServiceInstance() *EventsService {
	return eventsInstance
}

func (e *EventsService) AddEvent(ctx context.Context) {

}

func (e *EventsService) GetRaceEvents(ctx context.Context) {

}

func (e *EventsService) GetRaceEvent(ctx context.Context) {

}

func (e *EventsService) DeleteRaceEvent(ctx context.Context) {

}
