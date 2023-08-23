package v1

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type (
	ParticipantRequest struct {
		BibNumber string `json:"bib_number"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		BirthYear string `json:"birth_year"`
		Team      string `json:"team"`
		Gender    string `json:"gender"`
	}
)

func ParticipantsRoutes(handler *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.listParticipants)
	r.Post("/", handler.createParticipant)
	r.Get("/next_bib", handler.getNextBibNumber)
	r.Post("/csv", handler.importParticipantsCSV)
	r.Route("/{participantID}", func(r chi.Router) {
		r.Use(participantCtx)
		r.Get("/", handler.getParticipant)
	})

	return r
}

func (api *Handler) importParticipantsCSV(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		reader, err := r.MultipartReader()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

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

			var filename = part.FileName()

			log.Info().Msgf("%s - %d", filename, len(b))

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
					continue
				}

				// TODO Handle Names with more than 2 pieces
				log.Info().Str("raceID", raceID).Str("row", strings.Join(row, " - ")).Send()
				nameSplit := strings.Split(row[1], " ")
				log.Info().Int("split", len(nameSplit)).Interface("name", nameSplit).Send()
				firstName := nameSplit[0]
				lastName := ""
				if len(nameSplit) == 2 {
					lastName = nameSplit[1]
				}
				log.Info().Str("firstName", firstName).Str("lastName", lastName).Send()
				birthYear, _ := strconv.Atoi(row[3])
				gender := "Unknown"
				if row[4] == "M" {
					gender = "Male"
				} else if row[4] == "F" {
					gender = "Female"
				} else {
					log.Error().Msg("Skipping Row bad Gender")
					continue
				}
				pRow := services.ParticipantRow{
					RaceID:    raceID,
					BibNumber: row[0],
					FirstName: firstName,
					LastName:  lastName,
					Team:      row[2],
					BirthYear: birthYear,
					Gender:    gender,
				}

				err = services.AddParticipant(r.Context(), api.db, pRow)
				if err != nil {
					log.Error().Err(err).Str("value", strings.Join(row, ",")).Send()
				}
			}
		}
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (api *Handler) getNextBibNumber(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (api *Handler) listParticipants(w http.ResponseWriter, r *http.Request) {
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

	results, err := services.ListParticipants(r.Context(), api.db, limit, offset)
	if err != nil {
		log.Error().Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed request"))
		return
	}

	render.JSON(w, r, results)
}

func (api *Handler) getParticipant(w http.ResponseWriter, r *http.Request) {
}

func (api *Handler) createParticipant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := &ParticipantRequest{}
	err := render.DecodeJSON(r.Body, request)
	if err != nil {
		handleBadRequest(err, w, r)
		return
	}

	birthYear, err := strconv.Atoi(request.BirthYear)
	if err != nil {
		handleBadRequest(err, w, r)
		return
	}

	row := services.ParticipantRow{
		BibNumber: request.BibNumber,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		BirthYear: birthYear,
		Gender:    request.Gender,
		Team:      request.Team,
	}

	err = services.AddParticipant(ctx, api.db, row)
	if err != nil {
		handleBadRequest(err, w, r)
		return
	}

	render.Status(r, http.StatusCreated)
}

func participantCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ParticipantID := chi.URLParam(r, "participantID")
		if ParticipantID == "" {
			handleBadRequest(fmt.Errorf("unable to locate ParticipantID"), w, r)
			return
		}

		ctx := context.WithValue(r.Context(), util.ParticipantID, ParticipantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
