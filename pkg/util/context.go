package util

import "context"

type key int

const (
	OAuthTokenID  key = 0
	RaceID        key = 1
	EventID       key = 2
	GoogleClient  key = 3
	ParticipantID key = 4
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
