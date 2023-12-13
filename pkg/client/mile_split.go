package client

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Event struct {
	Name     string   `json:"name"`
	Results  []Result `json:"results"`
	Distance int      `json:"distance"`
}

type Result struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Place     int    `json:"place"`
	Team      string `json:"team"`
	Year      string `json:"year"`
	Time      int    `json:"time"`
}

var finishersRE = regexp.MustCompile(`\s?(\d+)\s+([\w-]+)\s([\w-]+)\s+([\w\s-]+)\s+((\d+:)?\d+.\d+)`)
var nonFinishersRE = regexp.MustCompile(`\s+(--)\s+([\w-]+)\s([\w-]+)\s+([\w\s-]+)\s+(DNS)`)
var headerRE = regexp.MustCompile(`\s*Athlete\s*Yr\s*Team\s*Time`)
var headerDistanceRE = regexp.MustCompile(`(\d+)\s(meter|kilometer|mile|km|m)`)

func GetEventInformation() Event {

	return Event{}
}

func getEventDistanceFromHeader(header string) int {
	matches := headerDistanceRE.FindStringSubmatch(header)
	if len(matches) != 3 {
		return 0
	}

	for _, match := range matches {
		fmt.Println(match)
	}

	raceDistance, err := strconv.ParseInt(matches[1], 0, 64)
	if err != nil {
		return 0
	}

	// return int(getDistanceToMetersConversionFactor(matches[2]) * raceDistance)

	return int(raceDistance)
}

func getDistanceInMeters(distance string, unit string) float64 {
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
		if distanceNum < 100 {
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
