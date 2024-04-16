package v1

import (
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/chriskuchin/roadrunner-results/pkg/db"
)

func HandleGetCurrentUser(db *db.DBLayer, auth *auth.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
