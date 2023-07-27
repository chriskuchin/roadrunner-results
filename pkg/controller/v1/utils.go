package v1

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

func handleBadRequest(err error, w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, map[string]string{
		"error": fmt.Sprintf("%v", err),
	})
}
