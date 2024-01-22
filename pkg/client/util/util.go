package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

var headerDistanceRE = regexp.MustCompile(`(?i)([\d,]+)\s?(meter|meters|kilometer|kilometers|mile|miles|km|m)`)
var headerGenderRE = regexp.MustCompile(`(?i)(boys|men|girls|women|male|female)`)
var maleGenderRE = regexp.MustCompile(`(?i)(boys|men|male)`)

// var femaleGenderRE = regexp.MustCompile(`(?i)(girls|women|female)`)

func GetEventDistanceFromHeader(header string) float64 {
	matches := headerDistanceRE.FindStringSubmatch(header)
	if len(matches) != 3 {
		return 0
	}

	raceDistance, err := strconv.ParseInt(strings.ReplaceAll(matches[1], ",", ""), 0, 64)
	if err != nil {
		return 0
	}

	return GetDistanceInMeters(fmt.Sprintf("%d", raceDistance), strings.ToLower(matches[2]))
}

func GetDistanceInMeters(distance string, unit string) float64 {
	distanceNum, err := strconv.ParseFloat(distance, 64)
	if err != nil {
		log.Warn().Err(err).Msg("Failed parsing mile split event distance")
		return 0
	}

	if unit == "meter" || unit == "meters" {
		return distanceNum
	} else if unit == "kilometer" || unit == "kilometers" || unit == "km" {
		return distanceNum * 1000
	} else if unit == "mile" || unit == "miles" {
		return distanceNum * 1609.344
	} else if unit == "m" {
		if distanceNum >= 50 {
			// meters
			return distanceNum
		} else {
			// miles
			return distanceNum * 1609.344
		}
	}

	log.Warn().Str("unit", unit).Msg("unknown unit")
	return 0
}

func GetGenderFromHeader(header string) string {
	matches := headerGenderRE.FindStringSubmatch(header)
	if len(matches) != 2 {
		return ""
	}

	if maleGenderRE.MatchString(strings.TrimSpace(matches[1])) {
		return "Male"
	} else {
		return "Female"
	}
}
