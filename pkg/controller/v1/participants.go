package v1

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func HandleParticipantsDistinctTeams() http.HandlerFunc {
	return apiutil.Unimplemented
}

func HandleParticipantsImportCSV(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if strings.HasPrefix(contentType, "multipart/form-data") {
			reader, err := r.MultipartReader()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			var rowsInserted int
			var failedRows int
			for {
				part, err := reader.NextPart()
				if err == io.EOF {
					break
				}

				b, err := io.ReadAll(part)
				if err != nil {
					log.Error().Err(err).Send()
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				reader := csv.NewReader(bytes.NewReader(b))

				// Iterate over the rows in the CSV file
				var raceID string = util.GetRaceIDFromContext(r.Context())
				var skipHeader bool = false
				for {
					// Read the next row
					row, err := reader.Read()
					if err == io.EOF {
						break
					}
					if !skipHeader {
						skipHeader = true
						continue
					}
					if row[1] == "" || len(row) != 5 {
						failedRows++
						log.Warn().Int("width", len(row)).Msg("not enough columns in the row")
						continue
					}

					firstName, lastName := util.SplitName(row[1])
					birthYear, _ := strconv.Atoi(row[3])
					gender := "Unknown"
					if row[4] == "M" || strings.ToLower(row[4]) == "male" {
						gender = "Male"
					} else if row[4] == "F" || strings.ToLower(row[4]) == "female" {
						gender = "Female"
					} else {
						log.Warn().Str("gender", gender).Msg("Skipping Row bad Gender")
						failedRows++
						continue
					}
					pRow := services.ParticipantRow{
						RaceID:    raceID,
						BibNumber: row[0],
						FirstName: firstName,
						LastName:  lastName,
						Team:      row[2],
						BirthYear: int64(birthYear),
						Gender:    gender,
					}

					err = services.AddParticipant(r.Context(), db, pRow)
					if err != nil {
						failedRows++
						log.Warn().Err(err).Str("value", strings.Join(row, ",")).Int("failed", failedRows).Int("success", rowsInserted).Send()
					} else {
						rowsInserted++
					}
				}
			}

			w.WriteHeader(http.StatusAccepted)
			render.JSON(w, r, map[string]int{
				"rows":   rowsInserted,
				"failed": failedRows,
			})
		}
	}
}

func HandleParticipantGetByBibNumber(db db.DB) http.HandlerFunc {
	type response struct {
		BibNumber string `json:"bib_number"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		BirthYear int64  `json:"birth_year"`
		Team      string `json:"team"`
		Gender    string `json:"gender"`
		Grade     int64  `json:"grade"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		participant, err := services.GetParticipantByBibNumber(ctx, db, util.GetRaceIDFromContext(ctx), util.GetBibNumberFromContext(ctx))
		if err != nil {
			if err == services.ErrorRowNotFound {
				apiutil.HandleNotFoundError(err, w, r)
				return
			}

			apiutil.HandleServerError(err, w, r)
			return
		}

		var result = response{
			BibNumber: participant.BibNumber,
			FirstName: participant.FirstName,
			LastName:  participant.LastName,
			Team:      participant.Team,
			BirthYear: participant.BirthYear,
			Gender:    participant.Gender,
		}

		if participant.Grade.Valid {
			result.Grade = participant.Grade.Int64
		}

		render.JSON(w, r, result)
	}
}

func HandleParticipantsNextBibNumber() http.HandlerFunc {
	return apiutil.Unimplemented
}

func HandleParticipantsList(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var err error
		var limit int = 10
		if r.URL.Query().Get("limit") != "" {
			limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Bad limit"))
				return
			}
		}

		var offset int = 0
		if r.URL.Query().Get("offset") != "" {
			offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Bad offset"))
				return
			}
		}

		var filters map[string]string = map[string]string{}
		if r.URL.Query().Get("gender") != "" {
			filters["gender"] = r.URL.Query().Get("gender")
		}

		if r.URL.Query().Get("team") != "" {
			filters["team"] = r.URL.Query().Get("team")
		}

		if r.URL.Query().Get("year") != "" {
			filters["year"] = r.URL.Query().Get("year")
		}

		if r.URL.Query().Get("name") != "" {
			filters["name"] = r.URL.Query().Get("name")
		}

		results, err := services.ListParticipants(ctx, db, util.GetRaceIDFromContext(ctx), limit, offset, filters)
		if err != nil {
			log.Error().Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed request"))
			return
		}

		render.JSON(w, r, results)
	}
}

func HandleParticipantGet() http.HandlerFunc {
	return apiutil.Unimplemented
}

func HandleParticipantUpdate(db db.DB) http.HandlerFunc {
	type ParticipantRequest struct {
		BibNumber int    `json:"bib_number"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		BirthYear int    `json:"birth_year"`
		Team      string `json:"team"`
		Gender    string `json:"gender"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		var payload ParticipantRequest
		err = json.Unmarshal(body, &payload)
		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		err = services.UpdateParticipant(ctx, db, util.GetRaceIDFromContext(ctx), util.GetParticipantIDFromContext(ctx), services.ParticipantRow{
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			BibNumber: fmt.Sprint(payload.BibNumber),
			Team:      payload.Team,
			Gender:    payload.Gender,
			BirthYear: int64(payload.BirthYear),
		})

		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleParticipantsCreate(db db.DB) http.HandlerFunc {
	type ParticipantRequest struct {
		BibNumber int    `json:"bib_number"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		BirthYear int    `json:"birth_year"`
		Team      string `json:"team"`
		Gender    string `json:"gender"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		request := &ParticipantRequest{}
		err := render.DecodeJSON(r.Body, request)
		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		row := services.ParticipantRow{
			BibNumber: fmt.Sprint(request.BibNumber),
			FirstName: request.FirstName,
			LastName:  request.LastName,
			BirthYear: int64(request.BirthYear),
			Gender:    request.Gender,
			Team:      request.Team,
			RaceID:    util.GetRaceIDFromContext(r.Context()),
		}

		err = services.AddParticipant(ctx, db, row)
		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		render.Status(r, http.StatusCreated)
	}
}
