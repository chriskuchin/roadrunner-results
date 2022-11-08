package services

import (
	"fmt"
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

func formatFinishingTime(ms int) string {
	current := ms

	min := current / msPerMin
	current = current % msPerMin

	sec := current / msPerSec
	current = current % msPerSec

	return fmt.Sprintf("%02d:%02d:%02d", min, sec, current/10)
}
