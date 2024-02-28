package middleware

import (
	"context"
	"net/http"
	"os"

	"firebase.google.com/go/v4/auth"
	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
)

func AuthenticationMiddleware(authClient *auth.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		token, token_present := os.LookupEnv("API_TOKEN")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("x-api-token")
			var tokResult *auth.Token
			var tokErr error
			if apiKey != "" {
				tokResult, tokErr = authClient.VerifyIDToken(r.Context(), apiKey)
				if tokErr == nil {
					r = r.WithContext(context.WithValue(r.Context(), util.UserToken, *tokResult))
				}
			}

			if r.Method != http.MethodGet && r.Method != http.MethodOptions {
				if (token_present && token == apiKey) || (tokErr == nil) {
					next.ServeHTTP(w, r)
					return
				} else {
					apiutil.HandleUnauthorized(w, r)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
