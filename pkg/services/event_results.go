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

func (ers *EventResultsService) GetEventResults(ctx context.Context) ([]ParticipantEventResult, error) {
	results, err := ers.eventResultsDao.GetEventResults(ctx)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	var participantResults []ParticipantEventResult = []ParticipantEventResult{}
	for _, result := range results {
		participantResults = append(participantResults, ParticipantEventResult(result))
	}

	return participantResults, nil
}

type ParticipantEventResult struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BibNumber string `json:"bib_number"`
	BirthYear int    `json:"birth_year"`
	Gender    string `json:"gender"`
	Result    int    `json:"result_ms"`
}
