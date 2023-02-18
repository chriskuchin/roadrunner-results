package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var spreadsheetIdRegex *regexp.Regexp = regexp.MustCompile(`https:\/\/docs.google.com\/spreadsheets\/d\/([\w\d_-]+)\/.*`)

type (
	Race struct {
		Name string
	}

	Event struct {
		Race        string
		Description string
		Distance    int
	}

	Participant struct {
		FirstName string
		LastName  string
		BirthYear string
		Gender    string
	}

	RaceParticipant struct {
		Participant string
		Race        string
		Event       string
		Team        string
		BibNumber   string
	}

	Result struct {
		Participant string
		Race        string
		Event       string
		Time        int
	}
)

var (
	race             Race
	events           []Event
	particpants      []Participant
	results          []Result
	raceParticipants []RaceParticipant
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// https://docs.google.com/spreadsheets/d/18Nj62AJHI-IbQaSn3dB_1TpGOmWFgQH8zMPKSkdC8fw/edit?usp=drivesdk
	spreadsheetId := "18Nj62AJHI-IbQaSn3dB_1TpGOmWFgQH8zMPKSkdC8fw"
	rslt, _ := srv.Spreadsheets.Get(spreadsheetId).Do()

	race = Race{
		Name: rslt.Properties.Title,
	}

	heats := []string{}
	for _, sheet := range rslt.Sheets {
		if isRaceHeatResultsTab(sheet.Properties.Title) {
			heats = append(heats, sheet.Properties.Title)
		}
	}

	for _, heat := range heats {
		events = append(events, Event{
			Race:        rslt.Properties.Title,
			Description: heat,
			Distance:    getHeatDistanceMeters(rslt.Properties.Title, heat),
		})

		readRange := fmt.Sprintf("%s!A:F", heat)
		resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).ValueRenderOption("FORMATTED_VALUE").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve data from sheet: %v", err)
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
					particpants = append(particpants, Participant{
						FirstName: strings.Split(row[1].(string), " ")[0],
						LastName:  strings.Join(strings.Split(row[1].(string), " ")[1:], " "),
						BirthYear: row[3].(string),
						Gender:    row[4].(string),
					})

					results = append(results, Result{
						Participant: row[1].(string),
						Race:        rslt.Properties.Title,
						Event:       heat,
						Time:        convertToMilliseconds(row[5].(string)),
					})

					raceParticipants = append(raceParticipants, RaceParticipant{
						Participant: row[1].(string),
						Race:        rslt.Properties.Title,
						Event:       heat,
						Team:        row[2].(string),
						BibNumber:   row[0].(string),
					})
				}
			}
			fmt.Println(len(particpants))
			fmt.Printf("Participant:     %+v\n", particpants[0])
			fmt.Printf("Result:          %+v\n", results[0])
			fmt.Printf("RaceParticipant: %+v\n", raceParticipants[0])
		}
	}

}

func extractSpreadsheetIDFromURL(url string) string {
	matches := spreadsheetIdRegex.FindStringSubmatch(url)
	if len(matches) != 2 {
		return ""
	}

	return matches[1]
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
