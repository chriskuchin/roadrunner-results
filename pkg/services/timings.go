package services

import (
	"strconv"
	"strings"
)

const (
	msPerMin = 60_000
	msPerSec = 1000
)

// 12:32:12
func parseFinishingTime(time string) (min int, sec int, hundredeths int) {
	timingPieces := strings.Split(time, ":")
	min, _ = strconv.Atoi(timingPieces[0])
	sec, _ = strconv.Atoi(timingPieces[1])
	hundredeths, _ = strconv.Atoi(timingPieces[2])

	return min, sec, hundredeths
}

func calculateMilisecondsFromTiming(min, sec, hund int) int {
	return msPerMin*min + msPerSec*sec + hund*10
}
