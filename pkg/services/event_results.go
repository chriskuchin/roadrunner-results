package services

import (
	"context"
	"time"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/rs/zerolog/log"
)

var eventResultsInstance *EventResultsService

type EventResultsService struct {
	eventResultsDao *db.EventResultsDao
	timerDao        *db.TimerDAO
	resultDao       *db.ResultsDAO
}

func NewEventResultsService() {
	if eventResultsInstance == nil {
		eventResultsInstance = &EventResultsService{
			eventResultsDao: db.NewEventResultsDAO(),
			timerDao:        db.NewTimerDAO(),
			resultDao:       db.NewResultsDAO(),
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

func (ers *EventResultsService) RecordResult(ctx context.Context, endTS int64) {
	start, timerID, err := ers.timerDao.GetActiveTimerStart(ctx)
	log.Info().Int64("start", start).Int64("end", endTS).Err(err).Send()

	timerStart := time.Unix(start/1000, (start%1000)*1000000)
	finishTime := time.Unix(endTS/1000, (endTS%1000)*1000000)

	log.Info().Msgf("%+v", finishTime.Sub(timerStart))

	log.Info().Str("timerID", timerID).Int64("result", endTS-start).Send()

	ers.timerDao.RecordTime(ctx, timerID, endTS-start)

}

type ParticipantEventResult struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BibNumber string `json:"bib_number"`
	BirthYear int    `json:"birth_year"`
	Gender    string `json:"gender"`
	Result    int    `json:"result_ms"`
}
