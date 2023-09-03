package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type VolunteerAuthPayload struct {
	Emails []VolunteerEmail `json:"emails"`
}

type VolunteerEmail struct {
	Email string `json:"email"`
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
	for _, email := range body.Emails {
		client, _ := api.app.Auth(ctx)
		user, err := client.GetUserByEmail(ctx, email.Email)
		if err != nil {
			failure = append(failure, user.Email)
			continue
		}

		if user.UID == util.GetCurrentUserID(ctx) {
			failure = append(failure, user.Email)
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
	w.WriteHeader(http.StatusNotImplemented)
}