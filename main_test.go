package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
)

var finishersRE = regexp.MustCompile(`\s?(\d+)\s+([\w-]+)\s([\w-]+)\s+([\w\s-]+)\s+((\d+:)?\d+.\d+)`)
var nonFinishersRE = regexp.MustCompile(`\s+(--)\s+([\w-]+)\s([\w-]+)\s+([\w\s-]+)\s+(DNS)`)
var headerRE = regexp.MustCompile(`\s*Athlete\s*Yr\s*Team\s*Time`)
var headerDistanceRE = regexp.MustCompile(`(\d+)\s(meter|kilometer|mile|km|m)`)

func Test_WebScrapingMileSplit(t *testing.T) {
	c := colly.NewCollector()

	c.OnHTML("div#meetResultsBody pre", func(e *colly.HTMLElement) {
		// fmt.Printf("\n\n%+v\n\n", e.Text)

		results := strings.Split(e.Text, "\n")

		for _, result := range results {
			if finishersRE.MatchString(result) {
				// fmt.Println("matched", result)
			} else if nonFinishersRE.MatchString(result) {
				// fmt.Printf("dns: \"%s\"\n", result)
			} else if !strings.HasPrefix(result, "===") && !headerRE.MatchString(result) && strings.TrimSpace(result) != "" {
				fmt.Printf("event: \"%s\"\n", result)
			}
			// \s(?<place>\d+)\s(?<first_name>\w+)\s(?<last_name>\w+)\s([\w\s-]+)(?=\s\d+\:\d+.\d+)\s(\d+:\d+.\d+)
		}
	})

	c.Visit("https://ut.milesplit.com/meets/511802-rep-distance-challenge-2023/results/876955/raw")
}

type ANetMeetInfo struct {
	MeetToken string `json:"jwtMeet"`
}

type ANetMeetResultsData struct {
	CurrentEventValid bool             `json:"currentEventValid"`
	Host              bool             `json:"host"`
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

func Test_WebScrapingAthleticNet(t *testing.T) {
	res, err := http.DefaultClient.Get("https://www.athletic.net/api/v1/Meet/GetMeetData?meetId=221697&sport=xc")
	if err != nil {
		log.Error().Err(err).Send()
		t.Fail()
	}

	meetInfo, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Send()
		t.Fail()
	}

	var meetObj ANetMeetInfo = ANetMeetInfo{}
	json.Unmarshal(meetInfo, &meetObj)

	log.Info().Interface("MeetInfo", meetObj).Send()

	body, _ := json.Marshal(map[string]string{"divId": "887176"})
	req, err := http.NewRequest(http.MethodPost, "https://www.athletic.net/api/v1/Meet/GetResultsData", bytes.NewBuffer(body))
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "ANETSettings=Team=45940&Sport=XC&guid=975e7ac4-974f-4653-9f0b-5ca79b32a30f")
	req.Header.Add("anettokens", meetObj.MeetToken)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	response, _ := io.ReadAll(res.Body)
	var raceResults ANetMeetResultsData
	json.Unmarshal(response, &raceResults)
	fmt.Printf("\n\n%+v\n\n", raceResults)
	//  results = "https://www.athletic.net/CrossCountry/meet/221697/results/887176")

}
