package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func HandleVolunteersCreate(db db.DB, app *firebase.App) http.HandlerFunc {
	type VolunteerEmail struct {
		Email string `json:"email"`
	}
	type VolunteerAuthPayload struct {
		Emails []VolunteerEmail `json:"emails"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload, err := io.ReadAll(r.Body)
		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		var body VolunteerAuthPayload
		err = json.Unmarshal(payload, &body)
		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		var success []string = []string{}
		var failure []string = []string{}
		client, err := app.Auth(ctx)
		if err != nil {
			log.Error().Err(err).Send()
			apiutil.HandleServerError(err, w, r)
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

			err = services.CreateRaceAuthorization(ctx, db, util.GetRaceIDFromContext(ctx), user.UID)
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
}

func HandleVolunteersList(db db.DB, app *firebase.App) http.HandlerFunc {
	type Volunteer struct {
		Email  string `json:"email"`
		UserID string `json:"userId"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ownerID, err := services.GetRaceOwnerID(ctx, db, util.GetRaceIDFromContext(ctx))
		if err != nil {
			log.Error().Err(err).Send()
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		if ownerID != util.GetCurrentUserID(ctx) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		userIds, err := services.ListRaceVolunteersUserIDs(ctx, db, util.GetRaceIDFromContext(ctx))
		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		client, err := app.Auth(ctx)
		if err != nil {
			// Need a server error
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		response := []Volunteer{}
		for _, userId := range userIds {
			ur, err := client.GetUser(ctx, userId)
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
}
