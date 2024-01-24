package v1

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type (
	ParticipantRequest struct {
		BibNumber int    `json:"bib_number"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		BirthYear int    `json:"birth_year"`
		Team      string `json:"team"`
		Gender    string `json:"gender"`
	}
)

func ParticipantsRoutes(handler *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.listParticipants)
	r.Post("/", handler.createParticipant)
	r.Get("/next_bib", handler.getNextBibNumber)
	r.Get("/teams", handler.getDistinctTeams)
	r.Post("/csv", handler.importParticipantsCSV)
	r.Route("/{participantID}", func(r chi.Router) {
		r.Use(participantCtx)
		r.Get("/", handler.getParticipant)
		r.Put("/", handler.updateParticipant)
	})

	return r
}

func (api *Handler) getDistinctTeams(w http.ResponseWriter, r *http.Request) {

}

func (api *Handler) importParticipantsCSV(w http.ResponseWriter, r *http.Request) {
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
					continue
				}

				firstName, lastName := util.SplitName(row[1])
				birthYear, _ := strconv.Atoi(row[3])
				gender := "Unknown"
				if row[4] == "M" {
					gender = "Male"
				} else if row[4] == "F" {
					gender = "Female"
				} else {
					log.Warn().Msg("Skipping Row bad Gender")
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

	results, err := services.ListParticipants(r.Context(), api.db, limit, offset, filters)
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

func (api *Handler) updateParticipant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(err, w, r)
		return
	}

	var payload ParticipantRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		handleBadRequest(err, w, r)
		return
	}

	err = services.UpdateParticipant(ctx, api.db, util.GetRaceIDFromContext(ctx), util.GetParticipantIDFromContext(ctx), services.ParticipantRow{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		BibNumber: fmt.Sprint(payload.BibNumber),
		Team:      payload.Team,
		Gender:    payload.Gender,
		BirthYear: payload.BirthYear,
	})

	if err != nil {
		handleBadRequest(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *Handler) createParticipant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := &ParticipantRequest{}
	err := render.DecodeJSON(r.Body, request)
	if err != nil {
		handleBadRequest(err, w, r)
		return
	}

	row := services.ParticipantRow{
		BibNumber: fmt.Sprint(request.BibNumber),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		BirthYear: request.BirthYear,
		Gender:    request.Gender,
		Team:      request.Team,
		RaceID:    util.GetRaceIDFromContext(r.Context()),
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
