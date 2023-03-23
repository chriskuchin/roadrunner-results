package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var timerServiceInstance *TimersService

type TimersService struct {
	dao *db.TimerDAO
}

type TimerResult struct {
	TimerID string `json:"id"`
	RaceID  string `json:"race_id"`
	EventID string `json:"event_id"`
	Start   int64  `json:"timer_start"`
}

func newTimersService() {
	if timerServiceInstance == nil {
		timerServiceInstance = &TimersService{
			dao: db.NewTimerDAO(),
		}
	}
}

func GetTimerServiceInstance() *TimersService {
	if timerServiceInstance == nil {
		newTimersService()
	}
	return timerServiceInstance
}

func (t *TimersService) StartTimer(ctx context.Context, start int64) {
	timerID := uuid.New().String()
	t.dao.StartTimer(ctx, timerID, start)
}

func (t *TimersService) ListTimers(ctx context.Context, limit, offset int) ([]TimerResult, error) {
	rows, err := t.dao.ListTimers(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	var results []TimerResult
	for _, row := range rows {

		results = append(results, TimerResult{
			TimerID: row.TimerID,
			RaceID:  row.RaceID,
			EventID: row.EventID,
			Start:   row.Start,
		})
	}

	return results, nil
}

func (t *TimersService) DeleteTimer(ctx context.Context) error {

	return t.dao.DeleteTimer(ctx)
}

func (t *TimersService) GetTimer(ctx context.Context) (*TimerResult, error) {
	result, err := t.dao.GetTimer(ctx)
	if err != nil {
		return nil, err
	}

	return &TimerResult{
		RaceID:  result.RaceID,
		EventID: result.EventID,
		TimerID: result.TimerID,
		Start:   result.Start,
	}, nil
}
