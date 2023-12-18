package client

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/gocolly/colly"
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
var resultExtractionRE = regexp.MustCompile(`\s(\d+)\s+(\w+)\s(\w+)\s+(\w+)\s([\w\s-]+)\s(\d+:\d+.\d+)`)

func GetEventInformation() Event {

	return Event{}
}

func processRawEventData(data string) []Event {
	var results []Event = []Event{}
	var currentEvent Event
	for _, result := range strings.Split(data, "\n") {
		if finishersRE.MatchString(result) {
			// fmt.Println("matched", result)
			getResultFromRow(result)
		} else if nonFinishersRE.MatchString(result) {
			// fmt.Printf("dns: \"%s\"\n", result)
		} else if !strings.HasPrefix(result, "===") && !headerRE.MatchString(result) && strings.TrimSpace(result) != "" {
			if currentEvent.Name != "" {
				results = append(results, currentEvent)
			}
			fmt.Printf("event: \"%s\"\n", result)
			currentEvent = Event{
				Name:     result,
				Results:  []Result{},
				Distance: getEventDistanceFromHeader(result),
			}
		}
		// \s(?<place>\d+)\s(?<first_name>\w+)\s(?<last_name>\w+)\s([\w\s-]+)(?=\s\d+\:\d+.\d+)\s(\d+:\d+.\d+)
	}

	fmt.Println(results)
	return results
}

func getResultFromRow(row string) (Result, error) {
	result := resultExtractionRE.FindAllStringSubmatch(row, -1)
	if len(result) > 0 && len(result[0]) == 7 {
		fmt.Println(strings.TrimSpace(result[0][6]))

		place, err := strconv.ParseInt(strings.TrimSpace(result[0][1]), 0, 62)
		if err != nil {
			return Result{}, fmt.Errorf("failed to parse place: %v", err)
		}

		return Result{
			Place:     int(place),
			FirstName: strings.TrimSpace(result[0][2]),
			LastName:  strings.TrimSpace(result[0][3]),
			Year:      strings.TrimSpace(result[0][4]),
			Team:      strings.TrimSpace(result[0][5]),
			Time:      util.ConvertFormatTimeToMilliseconds(strings.TrimSpace(result[0][6])),
		}, nil
	}

	return Result{}, fmt.Errorf("non result row failed to parse: %s", row)
}

func getRawEventData(eventURL string) string {
	var results string
	c := colly.NewCollector()

	c.OnHTML("div#meetResultsBody pre", func(e *colly.HTMLElement) {
		results = e.Text
	})
	c.Visit(eventURL)
	return results
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
