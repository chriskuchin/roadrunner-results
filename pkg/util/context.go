package util

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type key int

const (
	OAuthTokenID  key = 0
	RaceID        key = 1
	EventID       key = 2
	GoogleClient  key = 3
	ParticipantID key = 4
	TimerID       key = 5
	DB            key = 6
)

func getValueFromContext(ctx context.Context, k key) string {
	val := ctx.Value(k)
	result, ok := val.(string)
	if !ok {
		return ""
	}

	return result
}

func GetEventIDFromContext(ctx context.Context) string {
	return getValueFromContext(ctx, EventID)
}

func GetRaceIDFromContext(ctx context.Context) string {
	return getValueFromContext(ctx, RaceID)
}

func GetParticipantIDFromContext(ctx context.Context) string {
	return getValueFromContext(ctx, ParticipantID)
}

func GetTimerIDFromContext(ctx context.Context) string {
	return getValueFromContext(ctx, TimerID)
}

func SetDB(ctx context.Context, db *sqlx.DB) context.Context {
	return context.WithValue(ctx, DB, db)
}

func GetDB(ctx context.Context) *sqlx.DB {
	db := ctx.Value(DB)
	if db == nil {
		log.Fatal().Msg("db not found in context")
	}

	return db.(*sqlx.DB)
}
