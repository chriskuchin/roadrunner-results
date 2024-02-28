package middleware

import (
	"context"
	"fmt"
	"net/http"

	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi/v5"
)

func BibNumberCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bibNumber := chi.URLParam(r, "bib_number")
		if bibNumber == "" {
			apiutil.HandleBadRequest(fmt.Errorf("unable to locate bib_number"), w, r)
			return
		}

		ctx := context.WithValue(r.Context(), util.BibNumber, bibNumber)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
