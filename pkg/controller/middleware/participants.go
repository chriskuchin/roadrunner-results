package middleware

import (
	"context"
	"fmt"
	"net/http"

	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi/v5"
)

func ParticipantCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ParticipantID := chi.URLParam(r, "participantID")
		if ParticipantID == "" {
			apiutil.HandleBadRequest(fmt.Errorf("unable to locate ParticipantID"), w, r)
			return
		}

		ctx := context.WithValue(r.Context(), util.ParticipantID, ParticipantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
