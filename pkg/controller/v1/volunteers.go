package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type VolunteerAuthPayload struct {
	Email string
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
	json.Unmarshal(payload, &body)

	client, _ := api.app.Auth(ctx)
	user, _ := client.GetUserByEmail(ctx, body.Email)
	log.Info().Interface("user", user).Send()
	// services.CreateRaceAuthorization(ctx, api.db, util.GetRaceIDFromContext(ctx), user.UserInfo.UID)
	render.Status(r, http.StatusNotImplemented)
}

func (api *Handler) listVolunteers(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)
}
