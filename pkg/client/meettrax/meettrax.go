package meettrax

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/chriskuchin/roadrunner-results/pkg/client/model"
	"github.com/rs/zerolog/log"
)

var (
	meetIDRegex         = regexp.MustCompile(`https:\/\/meettrax.com\/meets\/(\d+).+`)
	inertiaVersionRegex = regexp.MustCompile(`&quot;version&quot;:&quot;(\w+\d+)&quot;`)
)

type MTTRXEventObject struct {
	Component string                `json:"componenet"`
	URL       string                `json:"url"`
	Version   string                `json:"version"`
	Props     MTTRXEventObjectProps `json:"props"`
}

type MTTRXEventObjectProps struct {
	VM MTTRXEventObjectVM `json:"vm"`
}

type MTTRXEventObjectVM struct {
	Events  MTTRXEventDataEvents  `json:"op_meet_events"`
	Results MTTRXEventResultsData `json:"op_meet_event_round_results"`
}

type MTTRXEventResultsData struct {
	Results MTTRXResultData `json:"results"`
}

type MTTRXResultData struct {
	Data []MTTRXRoundAthleteResultData `json:"data"`
}

type MTTRXMark struct {
	Type    string `json:"mark_type"`
	English string `json:"mark_english"`
	Metric  string `json:"mark_metric"`
}

type MTTRXRoundAthleteResultData struct {
	ID        int        `json:"id"`
	TeamName  string     `json:"team_name"`
	FirstName string     `json:"athlete_first_name"`
	LastName  string     `json:"athlete_last_name"`
	Grade     int        `json:"athlete_grade"`
	Place     int        `json:"meet_event_round_place"`
	Mark      *MTTRXMark `json:"mark"`
}

type MTTRXEventDataEvents struct {
	Data []MTTRXEventDataData `json:"data"`
}

type MTTRXEventDataData struct {
	ID          int    `json:"id"`
	DisplayName string `json:"display_name"`
	Gender      string `json:"gender"`
}

type MTTRXRawResults struct {
	EventID int
	Gender  string
	Name    string
	Result  []MTTRXRawAthleteResults
}

type MTTRXRawAthleteResults struct {
	TeamName  string
	FirstName string
	Lastname  string
	Grade     int
	Place     int
	Time      string
}

func GetEventInformation(eventURL string) []model.Event {
	meetId := getMeetIDFromURL(eventURL)
	rawEventData, err := getRawEventData(meetId)
	if err != nil {
		return nil
	}

	return processRawEventData(rawEventData)
}

func getMeetIDFromURL(eventURL string) string {
	found := meetIDRegex.FindAllStringSubmatch(eventURL, 1)

	if len(found[0]) == 2 {
		return found[0][1]
	}
	return ""
}

func processRawEventData([]MTTRXRawResults) []model.Event {
	return nil
}

func getRawEventData(meetID string) ([]MTTRXRawResults, error) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://meettrax.com/meets/%s/results/by-event", meetID), nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	found := inertiaVersionRegex.FindAllStringSubmatch(string(body), 1)

	versionID := found[0][1]

	req.Header.Add("X-Inertia", "1")
	req.Header.Add("X-Inertia-Version", versionID)

	finalResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	finalBody, _ := io.ReadAll(finalResponse.Body)

	var eventData MTTRXEventObject
	json.Unmarshal(finalBody, &eventData)

	var results []MTTRXRawResults = []MTTRXRawResults{}
	for _, event := range eventData.Props.VM.Events.Data {
		eventReq, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://meettrax.com/meets/%s/results/by-event/%d", meetID, event.ID), nil)
		eventReq.Header.Add("X-Inertia", "1")
		eventReq.Header.Add("X-Inertia-Version", versionID)

		eventRes, _ := http.DefaultClient.Do(eventReq)

		eventBody, _ := io.ReadAll(eventRes.Body)

		var eventResults MTTRXEventObject
		json.Unmarshal(eventBody, &eventResults)
		var rawEventResults MTTRXRawResults = MTTRXRawResults{
			EventID: event.ID,
			Name:    event.DisplayName,
			Gender:  event.Gender,
			Result:  []MTTRXRawAthleteResults{},
		}

		for _, result := range eventResults.Props.VM.Results.Results.Data {
			currentAthlete := MTTRXRawAthleteResults{
				TeamName:  result.TeamName,
				FirstName: result.FirstName,
				Lastname:  result.LastName,
				Place:     result.Place,
				Grade:     result.Grade,
			}

			if result.Mark != nil {
				currentAthlete.Time = result.Mark.English
				rawEventResults.Result = append(rawEventResults.Result, currentAthlete)
			}
		}

		results = append(results, rawEventResults)

	}

	return results, nil
}
