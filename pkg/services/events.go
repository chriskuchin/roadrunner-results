package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var eventsInstance *EventsService

type EventsService struct {
	eventsDao *db.EventDao
}

type EventObject struct {
	EventID     string `json:"eventId"`
	Description string `json:"description"`
	Distance    int    `json:"distance"`
}

func NewEventService() {
	if eventsInstance == nil {
		eventsInstance = &EventsService{
			eventsDao: db.NewEventDAO(),
		}
	}
}

func GetEventsServiceInstance() *EventsService {
	return eventsInstance
}

func (e *EventsService) AddEvent(ctx context.Context, raceID, description string, distance int) (string, error) {
	id := uuid.NewString()
	return id, e.eventsDao.AddEvent(ctx, raceID, id, description, distance)
}

func (e *EventsService) GetRaceEvents(ctx context.Context) ([]EventObject, error) {
	dbEvents, err := e.eventsDao.GetRaceEvents(ctx)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	results := []EventObject{}
	for _, event := range dbEvents {
		results = append(results, EventObject{
			EventID:     event.EventID,
			Description: event.Description,
			Distance:    event.Distance,
		})
	}

	return results, nil
}

func (e *EventsService) GetRaceEvent(ctx context.Context) {

}

func (e *EventsService) DeleteRaceEvent(ctx context.Context) {

}
