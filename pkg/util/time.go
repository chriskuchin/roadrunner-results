package util

import (
	"regexp"
	"strconv"
)

var parseTimeRE = regexp.MustCompile(`(\d+)`)

func ConvertFormatTimeToMilliseconds(formattedTime string) int {
	parsedTime := parseTimeRE.FindAllStringSubmatch(formattedTime, -1)

	var minutesIndex int = 0
	var secondsIndex int = 1
	var msIndex int = 2

	var err error
	var minutes int64 = 0
	if len(parsedTime) == 2 {
		secondsIndex = 0
		msIndex = 1
	} else if len(parsedTime) == 3 {
		minutes, err = strconv.ParseInt(parsedTime[minutesIndex][1], 0, 64)
		if err != nil {
			return 0
		}
	}

	seconds, err := strconv.ParseInt(parsedTime[secondsIndex][1], 0, 64)
	if err != nil {
		return 0
	}

	ms, err := strconv.ParseInt(parsedTime[msIndex][1], 0, 64)
	if err != nil {
		return 0
	}

	return int(minutes)*60000 + int(seconds)*1000 + int(ms)*10
	// let min = Math.floor(ms / 60000)
	// ms = ms % 60000
	// let sec = Math.floor(ms / 1000)
	// ms = Math.floor((ms % 1000) / 10)
	// return min + ":" + addLeadingZeros(sec) + "." + addLeadingZeros(ms)
}
