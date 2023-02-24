package google

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

func HandleOAuth2Creds(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken, err := r.Cookie("google-token")
		if err != nil {
			log.Error().Err(err).Send()

			currentURL := fmt.Sprintf("%s?%s", r.URL.Path, r.URL.RawQuery)
			url := GetOAuth2Config().AuthCodeURL(currentURL, oauth2.AccessTypeOffline)

			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		token, err := base64.StdEncoding.DecodeString(encodedToken.Value)
		if err != nil {
			log.Error().Err(err).Send()
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tok := &oauth2.Token{}
		err = json.Unmarshal([]byte(token), tok)
		if err != nil {
			log.Error().Err(err).Send()
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), util.OAuthTokenID, tok)
		ctx = context.WithValue(ctx, util.GoogleClient, getClient(tok))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
