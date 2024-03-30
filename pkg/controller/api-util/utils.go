package apiutil

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func HandleBadRequest(err error, w http.ResponseWriter, r *http.Request) {
	log.Warn().Err(err).Send()
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, map[string]string{
		"error": fmt.Sprintf("%v", err),
	})
}

func HandleServerError(err error, w http.ResponseWriter, r *http.Request) {
	log.Error().Err(err).Send()
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, map[string]string{
		"error": fmt.Sprintf("%v", err),
	})
}

func HandleServiceUnavailable(err error, w http.ResponseWriter, r *http.Request) {
	log.Error().Err(err).Send()
	render.Status(r, http.StatusServiceUnavailable)
	render.JSON(w, r, map[string]string{
		"error": fmt.Sprintf("%v", err),
	})
}

func HandleForbidden(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}

func HandleUnauthorized(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
}

func Unimplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
