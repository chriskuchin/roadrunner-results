package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type EventDao struct {
	db *sqlx.DB
}

func NewEventDAO(db *sqlx.DB) *EventDao {
	return &EventDao{
		db: db,
	}
}

func (e *EventDao) AddEvent(ctx context.Context) {

}

func (e *EventDao) GetRaceEvents(ctx context.Context) {

}

func (e *EventDao) DeleteEvent(ctx context.Context) {

}

func (e *EventDao) GetRaceEvent(ctx context.Context) {

}
