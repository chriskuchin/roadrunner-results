package middleware

import (
	"context"
	"net/http"

	"firebase.google.com/go/v4/auth"
	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
)

func AuthenticationMiddleware(verifyToken func(context.Context, string) (*auth.Token, error), allowedMethods []string, getenv func(string) string) func(http.Handler) http.Handler {
	var isAllowedMethod map[string]bool = map[string]bool{}
	for _, method := range allowedMethods {
		isAllowedMethod[method] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("x-api-token")
			var tokResult *auth.Token
			var tokErr error
			if apiKey != "" {
				tokResult, tokErr = verifyToken(r.Context(), apiKey)
				if tokErr == nil {
					r = r.WithContext(context.WithValue(r.Context(), util.UserToken, *tokResult))
				}
			}

			if isAllowedMethod[r.Method] || (tokErr == nil && tokResult != nil) {
				next.ServeHTTP(w, r)
				return
			}

			apiutil.HandleUnauthorized(w, r)
		})
	}
}

func UserAuthMiddleware(db *db.DBLayer, allowedMethods []string) func(http.Handler) http.Handler {
	var isAllowedMethod map[string]bool = map[string]bool{}
	for _, method := range allowedMethods {
		isAllowedMethod[method] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := isAllowedMethod[r.Method]; ok {
				next.ServeHTTP(w, r)
				return
			}

			raceID := util.GetRaceIDFromContext(r.Context())
			uid := util.GetCurrentUserID(r.Context())
			ownerID, err := services.GetRaceOwnerID(r.Context(), db, raceID)
			if err != nil {
				apiutil.HandleServiceUnavailable(err, w, r)
				return
			}

			if uid == ownerID {
				next.ServeHTTP(w, r)
				return
			}

			authorizedUsers, err := services.GetRaceAuthorizedUsers(r.Context(), db, raceID)
			if err != nil {
				apiutil.HandleServiceUnavailable(err, w, r)
				return
			}

			if authorizedUsers[uid] {
				next.ServeHTTP(w, r)
				return
			}

			apiutil.HandleUnauthorized(w, r)
		})
	}
}
