package services

import (
	"context"
	"encoding/json"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/rs/zerolog/log"
)

const (
	createRaceDivisionQuery string = `
		insert into race_divisions(
			race_id,
			display,
			filters
		)
		values (?,?,?)
	`

	listRaceDivisionsQuery = `
		select * from race_divisions
			where
				race_id = ?
			order by display
	`
)

type Filter struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

type DivisionRow struct {
	RaceID  string `db:"race_id"`
	Display string `db:"display"`
	Filters string `db:"filters"`
}

type DivisionResults struct {
	RaceID  string   `json:"race_id"`
	Display string   `json:"display"`
	Filters []Filter `json:"filters"`
}

func CreateDivision(ctx context.Context, db db.DB, raceID string, display string, filters []Filter) error {

	dbFilters, err := json.Marshal(filters)
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}
	_, err = db.ExecContext(ctx, createRaceDivisionQuery, raceID, display, string(dbFilters))
	if err != nil {
		return err
	}

	return nil
}

func ListRaceDivisions(ctx context.Context, db db.DB, raceID string) ([]DivisionResults, error) {
	divisionRows := []DivisionRow{}
	err := db.SelectContext(ctx, &divisionRows, listRaceDivisionsQuery, raceID)
	if err != nil {
		return nil, err
	}

	var divisions []DivisionResults = []DivisionResults{}
	for _, dr := range divisionRows {
		var filter []Filter

		json.Unmarshal([]byte(dr.Filters), &filter)
		divisions = append(divisions, DivisionResults{
			RaceID:  dr.RaceID,
			Display: dr.Display,
			Filters: filter,
		})
	}

	return divisions, nil
}
