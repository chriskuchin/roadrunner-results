package google

import (
	"context"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetOAuth2Config() *oauth2.Config {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Error().Msgf("Unable to read client secret file: %v", err)
	}

	conf, _ := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")

	return conf
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(token *oauth2.Token) *http.Client {
	return GetOAuth2Config().Client(context.Background(), token)
}
