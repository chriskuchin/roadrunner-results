package services

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func ImportFromSheet(ctx context.Context, sheetId string) {
	log.Info().Msg("Start Import")

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

	log.Info().Str("sheet", rslt.Properties.Title).Send()

	err = GetRaceServiceInstance().CreateRaceWithID(ctx, sheetId, rslt.Properties.Title)
	if err != nil {
		log.Error().Err(err).Send()
	}

	eventService := GetEventsServiceInstance()
	heats := []string{}
	for _, sheet := range rslt.Sheets {
		if isRaceHeatResultsTab(sheet.Properties.Title) {
			heats = append(heats, sheet.Properties.Title)
			err := eventService.AddEvent(ctx, sheetId, sheet.Properties.Title, getHeatDistanceMeters(sheet.Properties.Title))
			if err != nil {
				log.Error().Err(err).Send()
			}
		}
	}

	log.Info().Strs("heats", heats).Send()

	for _, heat := range heats {

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
					err = GetMembersServiceInstance().UpsertMember(ctx, strings.Split(row[1].(string), " ")[0], strings.Join(strings.Split(row[1].(string), " ")[1:], " "), row[4].(string), birthYear)
					if err != nil {
						log.Error().Err(err).Send()
					}
					// results = append(results, Result{
					// 	Participant: row[1].(string),
					// 	Race:        rslt.Properties.Title,
					// 	Event:       heat,
					// 	Time:        convertToMilliseconds(row[5].(string)),
					// })

					// raceParticipants = append(raceParticipants, RaceParticipant{
					// 	Participant: row[1].(string),
					// 	Race:        rslt.Properties.Title,
					// 	Event:       heat,
					// 	Team:        row[2].(string),
					// 	BibNumber:   row[0].(string),
					// })
				}
			}
			// fmt.Println(len(particpants))
			// fmt.Printf("Participant:     %+v\n", particpants[0])
			// fmt.Printf("Result:          %+v\n", results[0])
			// fmt.Printf("RaceParticipant: %+v\n", raceParticipants[0])
		}
	}

}

func convertToMilliseconds(timing string) int {
	split := strings.Split(timing, ":")
	seconds := strings.Split(split[1], ".")

	min, _ := strconv.Atoi(split[0])
	sec, _ := strconv.Atoi(seconds[0])
	tenth, _ := strconv.Atoi(seconds[1])

	result := min*60000 + sec*1000 + tenth*10

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
