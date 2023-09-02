package services

import (
	"context"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// {
// 	"actions": ["volunteer"],
// 	"events" : ["all"],
// }

// volunteer: no race modification
//            no delete
//            edit and create only

// owner: all permissions

// all: all events

// :id: specific event to manage

type PermissionPayload struct {
	Actions []string `json:"actions"`
	Events  []string `json:"events"`
}

const INSERT_RACE_AUTHORIZATION = `INSERT into race_authorization (user_id, race_id, permissions)
	VALUES(?, ?, json(?))`

func CreateRaceAuthorization(ctx context.Context, db *sqlx.DB, raceID, userID string) error {
	payload, _ := json.Marshal(PermissionPayload{
		Actions: []string{"volunteer"},
		Events:  []string{"all"},
	})

	result, err := db.ExecContext(ctx, INSERT_RACE_AUTHORIZATION, userID, raceID, string(payload))
	if err != nil {
		return err
	}

	log.Info().Interface("result", result).Send()
	return nil
}
