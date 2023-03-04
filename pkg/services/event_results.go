package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/rs/zerolog/log"
)

var eventResultsInstance *EventResultsService

type EventResultsService struct {
	eventResultsDao *db.EventResultsDao
}

func NewEventResultsService(dao *db.EventResultsDao) {
	if eventResultsInstance == nil {
		eventResultsInstance = &EventResultsService{
			eventResultsDao: dao,
		}
	}
}

func GetEventResultsInstance() *EventResultsService {
	return eventResultsInstance
}

func (ers *EventResultsService) GetEventResults(ctx context.Context) {
	log.Info().Msg("TESTTEST")
	err := ers.eventResultsDao.GetEventResults(ctx)
	log.Error().Err(err).Send()
}
