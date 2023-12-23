package mile_split

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/chriskuchin/roadrunner-results/pkg/client/model"
	client_util "github.com/chriskuchin/roadrunner-results/pkg/client/util"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/gocolly/colly"
)

var finishersRE = regexp.MustCompile(`\s?(\d+)\s+([\w-]+)\s([\w-]+)\s+([\w\s-]+)\s+((\d+:)?\d+.\d+)`)
var nonFinishersRE = regexp.MustCompile(`\s+(--)\s+([\w-]+)\s([\w-]+)\s+([\w\s-]+)\s+(DNS)`)
var headerRE = regexp.MustCompile(`\s*Athlete\s*Yr\s*Team\s*Time`)
var resultExtractionRE = regexp.MustCompile(`(\d+)\s+([\w-]+)\s([\w-]+)\s+(\w+)?\s([\w\s-]+)\s([\d+:\\.]+)`)

func GetEventInformation(eventURL string) []model.Event {

	rawEventData := getRawEventData(eventURL)

	return processRawEventData(rawEventData)
}

func processRawEventData(data string) []model.Event {
	var results []model.Event = []model.Event{}
	var currentEvent model.Event
	var currentGender string
	for _, result := range strings.Split(data, "\n") {
		if finishersRE.MatchString(result) {
			finisher, err := getResultFromRow(result, currentGender)
			if err != nil {
				continue
			}

			currentEvent.Results = append(currentEvent.Results, finisher)
		} else if nonFinishersRE.MatchString(result) {
			// fmt.Printf("dns: \"%s\"\n", result)
		} else if !strings.HasPrefix(result, "===") && !headerRE.MatchString(result) && strings.TrimSpace(result) != "" {
			if currentEvent.Name != "" {
				results = append(results, currentEvent)
			}
			currentEvent = model.Event{
				Name:     result,
				Results:  []model.Result{},
				Distance: client_util.GetEventDistanceFromHeader(result),
			}
			currentGender = client_util.GetGenderFromHeader(result)
		}
	}

	return results
}

func getResultFromRow(row, gender string) (model.Result, error) {
	result := resultExtractionRE.FindAllStringSubmatch(strings.TrimSpace(row), -1)
	if len(result) > 0 && len(result[0]) == 7 {
		place, err := strconv.ParseInt(strings.TrimSpace(result[0][1]), 0, 62)
		if err != nil {
			return model.Result{}, fmt.Errorf("failed to parse place: %v", err)
		}

		return model.Result{
			Place:     int(place),
			FirstName: strings.TrimSpace(result[0][2]),
			LastName:  strings.TrimSpace(result[0][3]),
			Year:      strings.TrimSpace(result[0][4]),
			Team:      strings.TrimSpace(result[0][5]),
			Gender:    gender,
			Time:      util.ConvertFormatTimeToMilliseconds(strings.TrimSpace(result[0][6])),
		}, nil
	}

	return model.Result{}, fmt.Errorf("non result row failed to parse: %s", row)
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
