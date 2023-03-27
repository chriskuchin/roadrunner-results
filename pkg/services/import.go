package services

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func ImportFromSheet(ctx context.Context, db *sqlx.DB, sheetId string) {
	sheets, err := sheets.NewService(ctx, option.WithHTTPClient(ctx.Value(util.GoogleClient).(*http.Client)))
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	rslt, err := sheets.Spreadsheets.Get(sheetId).Do()
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	err = CreateRaceWithID(ctx, db, sheetId, rslt.Properties.Title)
	if err != nil {
		log.Error().Err(err).Send()
	}

	heats := []string{}
	events := []string{}
	for _, sheet := range rslt.Sheets {
		if isRaceHeatResultsTab(sheet.Properties.Title) {
			heats = append(heats, sheet.Properties.Title)
			eventID, err := AddEvent(ctx, db, sheetId, sheet.Properties.Title, getHeatDistanceMeters(sheet.Properties.Title))
			if err != nil {
				log.Error().Err(err).Send()
			}
			events = append(events, eventID)
		}
	}

	for heatID, heat := range heats {

		readRange := fmt.Sprintf("%s!A:F", heat)
		resp, err := sheets.Spreadsheets.Values.Get(sheetId, readRange).ValueRenderOption("FORMATTED_VALUE").Do()
		if err != nil {
			log.Fatal().Msgf("Unable to retrieve data from sheet: %v", err)
		}

		if len(resp.Values) == 0 {
			fmt.Println("No data found.")
		} else {
			for idx, row := range resp.Values {
				if idx == 0 || len(row) <= 1 {
					continue
				} else {
					if row[0].(string) == "" {
						break
					}
					birthYear, _ := strconv.Atoi(row[3].(string))
					AddParticipant(ctx, db, Participant{
						RaceID: sheetId,
						// EventID:   events[heatID],
						FirstName: strings.Split(row[1].(string), " ")[0],
						LastName:  strings.Join(strings.Split(row[1].(string), " ")[1:], " "),
						BirthYear: birthYear,
						Gender:    row[4].(string),
						Team:      row[2].(string),
						BibNumber: row[0].(string),
					})
					if err != nil {
						log.Error().Err(err).Send()
					}

					err = RecordEventParticipation(ctx, db, sheetId, events[heatID], row[0].(string))

					if err != nil {
						log.Error().Err(err).Send()
					}

					err = InsertResults(ctx, db, sheetId, events[heatID], row[0].(string), row[5].(string))
					if err != nil {
						log.Error().Err(err).Send()
					}
				}
			}
		}
	}

}

func convertToMilliseconds(timing string) int64 {
	split := strings.Split(timing, ":")
	seconds := strings.Split(split[1], ".")

	min, _ := strconv.ParseInt(split[0], 10, 64)
	sec, _ := strconv.ParseInt(seconds[0], 10, 64)
	tenth, _ := strconv.ParseInt(seconds[1], 10, 64)

	var result int64 = min*60000 + sec*1000 + tenth*10

	return result
}

var raceID *regexp.Regexp = regexp.MustCompile(`Race|(\d+)([mMkK])`)
var raceDistance *regexp.Regexp = regexp.MustCompile(`(\d+)([mMkK])`)

func isRaceHeatResultsTab(tab string) bool {
	return raceID.MatchString(tab)
}

func getHeatDistanceMeters(description ...string) int {
	for _, race := range description {
		matches := raceDistance.FindStringSubmatch(race)

		if len(matches) != 3 {
			return 0
		}

		if strings.ToLower(matches[2]) == "k" {
			dist, _ := strconv.Atoi(matches[1])

			return dist * 1000
		} else if strings.ToLower(matches[2]) == "m" {
			dist, _ := strconv.Atoi(matches[1])
			return dist
		}
	}
	return 0
}
