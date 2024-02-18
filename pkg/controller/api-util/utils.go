package apiutil

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

func HandleBadRequest(err error, w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, map[string]string{
		"error": fmt.Sprintf("%v", err),
	})
}

func HandleServerError(err error, w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, map[string]string{
		"error": fmt.Sprintf("%v", err),
	})
}

func Unimplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
