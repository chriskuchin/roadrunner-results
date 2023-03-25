package services

import (
	"context"

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
		participantResults = append(participantResults, ParticipantEventResult{
			BibNumber: result.BibNumber,
			Result:    result.Result,
			FirstName: result.FirstName.String,
			LastName:  result.LastName.String,
			BirthYear: int(result.BirthYear.Int32),
			Gender:    result.Gender.String,
		})
	}

	return participantResults, nil
}

func (ers *EventResultsService) RecordTimerResult(ctx context.Context, endTS int64) error {
	start, _, err := ers.timerDao.GetActiveTimerStart(ctx)
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}

	return ers.resultDao.InsertPartialResult(ctx, endTS-start)
}

func (ers *EventResultsService) RecordFinisherResult(ctx context.Context, bib string) error {
	return ers.eventResultsDao.RecordFinisher(ctx, bib)
}

type ParticipantEventResult struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	BibNumber string `json:"bib_number"`
	BirthYear int    `json:"birth_year,omitempty"`
	Gender    string `json:"gender,omitempty"`
	Result    int    `json:"result_ms"`
}
