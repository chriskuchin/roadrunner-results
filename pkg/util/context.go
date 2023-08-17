package util

import (
	"context"

	"firebase.google.com/go/v4/auth"
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
	HeatID        key = 7
	PhotoKey      key = 8
	UserToken     key = 9
)

func GetCurrentUserToken(ctx context.Context) auth.Token {
	rawToken := ctx.Value(UserToken)
	if rawToken != nil {
		token, ok := rawToken.(auth.Token)
		if ok {
			return token
		}
	}
	return auth.Token{}
}

func GetCurrentUserID(ctx context.Context) string {
	return GetCurrentUserToken(ctx).UID
}

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

func GetPhotoKeyFromContext(ctx context.Context) string {
	return getValueFromContext(ctx, PhotoKey)
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
