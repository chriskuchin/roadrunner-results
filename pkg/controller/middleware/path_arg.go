package middleware

import (
	"context"
	"fmt"
	"net/http"

	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi/v5"
)

func PathArgCtx(argName string, key util.Key) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			val := chi.URLParam(r, argName)
			if val == "" {
				apiutil.HandleBadRequest(fmt.Errorf("unable to locate %s", argName), w, r)
				return
			}

			ctx := context.WithValue(r.Context(), key, val)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
