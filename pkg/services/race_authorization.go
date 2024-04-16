package services

import (
	"context"
	"encoding/json"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
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

func CreateRaceAuthorization(ctx context.Context, db db.DB, raceID, userID string) error {
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

const LOOKUP_RACE_AUTHORIZATION = `SELECT user_id FROM race_authorization WHERE race_id = ?`

func GetRaceAuthorizedUsers(ctx context.Context, db db.DB, raceID string) (map[string]bool, error) {
	rows, err := db.QueryContext(ctx, LOOKUP_RACE_AUTHORIZATION, raceID)
	if err != nil {
		return nil, err
	}

	var authorizedUser map[string]bool = map[string]bool{}
	for rows.Next() {
		var user_id string
		rows.Scan(&user_id)

		authorizedUser[user_id] = true
	}

	return authorizedUser, nil
}

const LOOKUP_EVENT_RACE_AUTHORIZATION = `SELECT user_id FROM race_authorization as ra JOIN json_each(race_authorization.permissions, "$.events") as events WHERE ra.race_id = ? AND (events.value = "all" OR events.value = ?);`

func GetEventAuthorizedUsers(ctx context.Context, db db.DB, raceID, eventID string) (map[string]bool, error) {
	rows, err := db.QueryContext(ctx, LOOKUP_EVENT_RACE_AUTHORIZATION, raceID, eventID)
	if err != nil {
		return nil, err
	}

	var authorizedUsers map[string]bool = map[string]bool{}
	for rows.Next() {
		var user_id string
		rows.Scan(&user_id)

		authorizedUsers[user_id] = true
	}
	return authorizedUsers, nil
}

const LIST_VOLUNTEER_USER_IDS = `SELECT user_id FROM race_authorization WHERE race_id = ?`

func ListRaceVolunteersUserIDs(ctx context.Context, db db.DB, raceID string) ([]string, error) {
	rows, err := db.QueryContext(ctx, LIST_VOLUNTEER_USER_IDS, raceID)
	if err != nil {
		return nil, err
	}

	volunteers := []string{}
	for rows.Next() {
		var user_id string
		rows.Scan(&user_id)

		volunteers = append(volunteers, user_id)
	}
	return volunteers, nil
}
