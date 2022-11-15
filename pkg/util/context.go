package util

import "context"

type key string

const (
	RaceID  key = "raceID"
	EventID key = "eventID"
)

func GetRaceIDFromContext(ctx context.Context) string {
	raceID := ctx.Value(RaceID)
	result, ok := raceID.(string)
	if !ok {
		return ""
	}

	return result
}
