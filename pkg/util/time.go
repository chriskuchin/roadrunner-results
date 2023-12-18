package util

import (
	"regexp"
	"strconv"
)

var parseTimeRE = regexp.MustCompile(`(\d+):(\d+).(\d+)`)

func ConvertFormatTimeToMilliseconds(formattedTime string) int {
	parsedTime := parseTimeRE.FindAllStringSubmatch(formattedTime, -1)
	minutes, err := strconv.ParseInt(parsedTime[0][1], 0, 64)
	if err != nil {
		return 0
	}

	seconds, err := strconv.ParseInt(parsedTime[0][2], 0, 64)
	if err != nil {
		return 0
	}

	ms, err := strconv.ParseInt(parsedTime[0][3], 0, 64)
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
