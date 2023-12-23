package athletic_net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/chriskuchin/roadrunner-results/pkg/client/model"
	"github.com/chriskuchin/roadrunner-results/pkg/client/util"
	"github.com/rs/zerolog/log"
)

type ANetMeetInfo struct {
	MeetToken string          `json:"jwtMeet"`
	Divisions []ANetDivisions `json:"xcDivisions"`
}

type ANetDivisions struct {
	Name       string `json:"DivName"`
	DivisionID int64  `json:"IDMeetDiv"`
}

type ANetMeetResultsData struct {
	CurrentEventValid bool `json:"currentEventValid"`
	Host              bool `json:"host"`
	Name              string
	Results           []ANetMeetResult `json:"results"`
}

type ANetMeetResult struct {
	ID          int64   `json:"IDResult"`
	Result      string  `json:"Result"`
	SortValue   float32 `json:"SortValue"`
	Place       int     `json:"Place"`
	AthleteName string  `json:"AthleteName"`
	TeamName    string  `json:"TeamName"`
	Official    int     `json:"Official"`
	AthleteID   int64   `json:"AthleteID"`
	FirstName   string  `json:"FirstName"`
	LastName    string  `json:"LastName"`
	Grade       string  `json:"Grade"`
	AgeGrade    string  `json:"AgeGrade"`
	Gender      string  `json:"Gender"`
	TeamID      int64   `json:"TeamID"`
	SchoolName  string  `json:"Schoolname"`
	TeamCode    string  `json:"TeamCode"`
}

var infoURLRE = regexp.MustCompile(`https://(www.)?athletic.net/(\w+)/meet/(\w+\d+)/info`)

func GetEventInformation(eventURL string) []model.Event {

	meetId, meetType := getEventIDAndTypeFromURL(eventURL)
	rawEventData, err := getRawEventData(meetId, meetType)
	if err != nil {
		return nil
	}

	return processRawEventData(rawEventData)
}

func processRawEventData(data []ANetMeetResultsData) []model.Event {
	var events []model.Event = []model.Event{}
	for _, event := range data {
		var currentEvent model.Event = model.Event{
			Name:     event.Name,
			Results:  []model.Result{},
			Distance: util.GetEventDistanceFromHeader(event.Name),
		}
		for _, result := range event.Results {
			var gender string = "Male"
			if result.Gender == "F" {
				gender = "Female"
			}
			currentEvent.Results = append(currentEvent.Results, model.Result{
				Place:     result.Place,
				FirstName: result.FirstName,
				LastName:  result.LastName,
				Time:      int(result.SortValue * 1000),
				Team:      result.TeamName,
				Year:      result.AgeGrade,
				Gender:    gender,
			})
		}
		events = append(events, currentEvent)
	}

	return events
}

func getRawEventData(meetID, meetType string) (data []ANetMeetResultsData, err error) {
	res, err := http.DefaultClient.Get(fmt.Sprintf("https://www.athletic.net/api/v1/Meet/GetMeetData?meetId=%s&sport=%s", meetID, meetType))
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	meetInfo, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	var meetObj ANetMeetInfo = ANetMeetInfo{}
	json.Unmarshal(meetInfo, &meetObj)

	data = []ANetMeetResultsData{}
	for _, division := range meetObj.Divisions {
		body, _ := json.Marshal(map[string]int64{"divId": division.DivisionID})
		log.Info().Str("body", string(body)).Send()
		req, err := http.NewRequest(http.MethodPost, "https://www.athletic.net/api/v1/Meet/GetResultsData", bytes.NewBuffer(body))
		if err != nil {
			log.Error().Err(err).Send()
			continue
		}
		req.Header.Add("content-type", "application/json")
		req.Header.Add("anettokens", meetObj.MeetToken)

		res, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Error().Err(err).Send()
			continue
		}

		response, _ := io.ReadAll(res.Body)
		var raceResults ANetMeetResultsData
		err = json.Unmarshal(response, &raceResults)
		if err != nil {
			log.Error().Err(err).Send()
			continue
		}

		raceResults.Name = division.Name

		data = append(data, raceResults)
	}

	payload, _ := json.Marshal(data)
	os.WriteFile("testdata/bob_firman_1_2023", payload, 0777)

	return data, nil
}

func getEventIDAndTypeFromURL(infoPageURL string) (meetID string, eventType string) {
	pieces := infoURLRE.FindAllStringSubmatch(infoPageURL, -1)

	meetID = pieces[0][3]

	if pieces[0][2] == "CrossCountry" {
		eventType = "xc"
	} else {
		eventType = "unknown"
	}

	return
}
