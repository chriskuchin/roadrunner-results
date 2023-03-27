package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func ParticipantsRoutes(handler *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.listParticipants)
	r.Route("/{participantID}", func(r chi.Router) {
		r.Use(participantCtx)
		r.Get("/", handler.getParticipant)
	})

	return r
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

func participantCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ParticipantID := chi.URLParam(r, "ParticipantID")
		if ParticipantID == "" {
			log.Error().Msg("Unable to locate ParticipantID")
		}

		ctx := context.WithValue(r.Context(), util.ParticipantID, ParticipantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
