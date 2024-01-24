package util

import (
	"regexp"
	"strconv"

	"github.com/rs/zerolog/log"
)

var parseTimeRE = regexp.MustCompile(`(\d+)`)

func ConvertFormatTimeToMilliseconds(formattedTime string) int {
	parsedTime := parseTimeRE.FindAllStringSubmatch(formattedTime, -1)

	var minutesIndex int = 0
	var secondsIndex int = 1
	var msIndex int = 2

	var err error
	var minutes int = 0
	if len(parsedTime) == 2 {
		secondsIndex = 0
		msIndex = 1
	} else if len(parsedTime) == 3 {
		minutes, err = strconv.Atoi(parsedTime[minutesIndex][1])
		if err != nil {
			log.Error().Str("time_input", formattedTime).Err(err).Send()
			return 0
		}
	}

	seconds, err := strconv.Atoi(parsedTime[secondsIndex][1])
	if err != nil {
		log.Error().Str("time_input", formattedTime).Err(err).Send()
		return 0
	}

	ms, err := strconv.Atoi(parsedTime[msIndex][1])
	if err != nil {
		log.Error().Str("time_input", formattedTime).Err(err).Send()
		return 0
	}

	return minutes*60000 + seconds*1000 + ms*10
}
