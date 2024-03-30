package v1

import (
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/jmoiron/sqlx"
)

func HandleGetCurrentUser(db *sqlx.DB, auth *auth.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
