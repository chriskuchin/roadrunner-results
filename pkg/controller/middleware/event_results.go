package middleware

import (
	"context"
	"fmt"
	"net/http"

	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi/v5"
)

func ResultCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resultID := chi.URLParam(r, "resultID")
		if resultID == "" {
			apiutil.HandleBadRequest(fmt.Errorf("unable to locate resultID"), w, r)
			return
		}

		ctx := context.WithValue(r.Context(), util.ResultID, resultID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
