package v1

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/chriskuchin/roadrunner-results/pkg/google"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

type googleResources struct{}

func (rs googleResources) Routes() chi.Router {
	r := chi.NewRouter()
	r.HandleFunc("/oauth2/callback", handleOAuth2Callback)
	return r
}

func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	token, err := google.GetOAuth2Config().Exchange(context.TODO(), code)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	marshaledToken, err := json.Marshal(token)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	cookieValue := base64.StdEncoding.EncodeToString(marshaledToken)
	http.SetCookie(w, &http.Cookie{Name: "google-token", Value: cookieValue, Expires: time.Now().Add(1 * time.Hour), Path: "/"})
	if state != "" {
		http.Redirect(w, r, state, http.StatusFound)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
