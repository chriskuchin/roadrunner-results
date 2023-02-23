package v1

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleResources struct{}

func (rs googleResources) Routes() chi.Router {
	r := chi.NewRouter()
	r.HandleFunc("/oauth2/callback", handleOAuth2Callback)
	return r
}

func getOAuth2Config() *oauth2.Config {
	// conf := &oauth2.Config{
	// 	ClientID:     "",
	// 	ClientSecret: "",
	// 	Scopes:       []string{"https://www.googleapis.com/auth/spreadsheets.readonly"},
	// 	Endpoint: oauth2.Endpoint{
	// 		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
	// 		TokenURL: "https://oauth2.googleapis.com/token",
	// 	},
	// 	RedirectURL: "https://debug.home.cksuperman.com/api/v1/google/oauth2/callback",
	// }

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Error().Msgf("Unable to read client secret file: %v", err)
	}

	conf, _ := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")

	return conf
}

func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	token, err := getOAuth2Config().Exchange(context.TODO(), code)
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
	http.SetCookie(w, &http.Cookie{Name: "oauth-token", Value: cookieValue, Expires: time.Now().Add(1 * time.Hour), Path: "/"})
	if state != "" {
		http.Redirect(w, r, state, http.StatusFound)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
