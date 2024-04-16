package util

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

type Key int

const (
	OAuthTokenID  Key = 0
	RaceID        Key = 1
	EventID       Key = 2
	GoogleClient  Key = 3
	ParticipantID Key = 4
	TimerID       Key = 5
	DB            Key = 6
	HeatID        Key = 7
	PhotoKey      Key = 8
	UserToken     Key = 9
	ResultID      Key = 10
	BibNumber     Key = 11
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

func getValueFromContext(ctx context.Context, k Key) string {
	val := ctx.Value(k)
	result, ok := val.(string)
	if !ok {
		return ""
	}

	return result
}

func GetBibNumberFromContext(ctx context.Context) string {
	return getValueFromContext(ctx, BibNumber)
}

func GetResultIDFromContext(ctx context.Context) string {
	return getValueFromContext(ctx, ResultID)
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
