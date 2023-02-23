package util

import "context"

type key int

const (
	OAuthTokenID key = 0
	RaceID       key = 1
	EventID      key = 2
)

func GetRaceIDFromContext(ctx context.Context) string {
	raceID := ctx.Value(RaceID)
	result, ok := raceID.(string)
	if !ok {
		return ""
	}

	return result
}
