package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type VolunteerAuthPayload struct {
	Emails []VolunteerEmail `json:"emails"`
}

type VolunteerEmail struct {
	Email string `json:"email"`
}

type Volunteer struct {
	Email  string `json:"email"`
	UserID string `json:"userId"`
}

func VolunteerRoutes(handler *Handler) chi.Router {
	r := chi.NewRouter()

	r.Put("/", handler.authorizeVolunteer)
	r.Get("/", handler.listVolunteers)

	return r
}

func (api *Handler) authorizeVolunteer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(err, w, r)
		return
	}

	var body VolunteerAuthPayload
	err = json.Unmarshal(payload, &body)
	if err != nil {
		handleBadRequest(err, w, r)
		return
	}

	var success []string = []string{}
	var failure []string = []string{}
	client, err := api.app.Auth(ctx)
	if err != nil {
		log.Error().Err(err).Send()
		handleServerError(err, w, r)
		return
	}

	for _, email := range body.Emails {
		user, err := client.GetUserByEmail(ctx, strings.TrimSpace(email.Email))
		if err != nil {
			log.Error().Err(err).Send()
			failure = append(failure, email.Email)
			continue
		}

		if user.UID == util.GetCurrentUserID(ctx) {
			failure = append(failure, email.Email)
			continue
		}

		err = services.CreateRaceAuthorization(ctx, api.db, util.GetRaceIDFromContext(ctx), user.UID)
		if err != nil {
			failure = append(failure, user.Email)
		} else {
			success = append(success, user.Email)
		}
	}
	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, map[string][]string{
		"success": success,
		"failure": failure,
	})
}

func (api *Handler) listVolunteers(w http.ResponseWriter, r *http.Request) {
	ownerID, err := services.GetRaceOwnerID(r.Context(), api.db, util.GetRaceIDFromContext(r.Context()))
	if err != nil {
		log.Error().Err(err).Send()
		handleBadRequest(err, w, r)
		return
	}

	if ownerID != util.GetCurrentUserID(r.Context()) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userIds, err := services.ListRaceVolunteersUserIDs(r.Context(), api.db, util.GetRaceIDFromContext(r.Context()))
	if err != nil {
		handleBadRequest(err, w, r)
		return
	}

	client, err := api.app.Auth(r.Context())
	if err != nil {
		// Need a server error
		handleBadRequest(err, w, r)
		return
	}

	response := []Volunteer{}
	for _, userId := range userIds {
		ur, err := client.GetUser(r.Context(), userId)
		if err != nil {
			continue
		}

		response = append(response, Volunteer{
			Email:  ur.Email,
			UserID: ur.UserInfo.UID,
		})
	}

	render.JSON(w, r, response)
}
