package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
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
	w.WriteHeader(http.StatusNotImplemented)
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
