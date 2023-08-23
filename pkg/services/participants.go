package services

import (
	"context"
	"fmt"

	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const (
	addParticipantQuery = `
		insert into participants (
			race_id,
			bib_number,
			first_name,
			last_name,
			gender,
			team,
			birth_year
		)
		VALUES (
			:race_id,
			:bib_number,
			:first_name,
			:last_name,
			:gender,
			:team,
			:birth_year
		)
	`
	selectRaceParticipantCountQuery = `
		select count(1) from participants
			where
				race_id = ?
	`

	selectRaceParticipantGenderCountQuery = `
		select gender, count(1) from participants
			where race_id = ?
			GROUP BY gender
	`

	selectRaceParticipantBirthYearCountQuery = `
		select birth_year, gender, count(1) from participants
			where race_id = ?
			GROUP BY birth_year, gender
			ORDER BY birth_year
	`
)

type (
	ParticipantRow struct {
		RaceID    string `db:"race_id"`
		BibNumber string `db:"bib_number"`
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		BirthYear int    `db:"birth_year"`
		Gender    string `db:"gender"`
		Team      string `db:"team"`
	}

	ParticipantObject struct {
		BibNumber string `json:"bibNumber"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		BirthYear int    `json:"birthYear"`
		Gender    string `json:"gender"`
		Team      string `json:"team"`
	}
)

func AddParticipant(ctx context.Context, db *sqlx.DB, participant ParticipantRow) error {
	_, err := db.NamedExec(addParticipantQuery, participant)
	return err
}

const listParticipantsQuery = `
select * from participants
	where
		race_id = ?
`

func ListParticipants(ctx context.Context, db *sqlx.DB, limit, offset int, filters map[string]string) ([]ParticipantObject, error) {
	var query string = listParticipantsQuery

	var values []interface{} = []interface{}{}
	values = append(values, util.GetRaceIDFromContext(ctx))
	for k, v := range filters {
		if k == "team" || k == "gender" {
			query = fmt.Sprintf("%s AND %s = ?", query, k)
			values = append(values, v)
		} else if k == "year" {
			query = fmt.Sprintf("%s AND birth_year = ?", query)
			values = append(values, v)
		} else if k == "bib" {
			query = fmt.Sprintf("%s AND bib_number = ?", query)
			values = append(values, v)
		} else if k == "name" {
			query = fmt.Sprintf("%s AND (first_name LIKE ? OR last_name LIKE ?)", query)
			values = append(values, fmt.Sprintf("%%%s", v))
			values = append(values, fmt.Sprintf("%%%s", v))
		}
	}

	query = fmt.Sprintf("%s limit ? offset ?", query)
	values = append(values, limit)
	values = append(values, offset)

	log.Info().Str("query", query).Interface("filters", filters).Send()

	var dbResults []ParticipantRow
	err := db.Select(&dbResults, query, values...)
	if err != nil {
		return nil, err
	}

	results := []ParticipantObject{}
	for _, participant := range dbResults {
		results = append(results, ParticipantObject{
			BibNumber: participant.BibNumber,
			FirstName: participant.FirstName,
			LastName:  participant.LastName,
			BirthYear: participant.BirthYear,
			Gender:    participant.Gender,
			Team:      participant.Team,
		})
	}

	return results, nil
}

func GetRaceGenderCount(ctx context.Context, db *sqlx.DB) (female int, male int, err error) {
	rows, err := db.Query(selectRaceParticipantGenderCountQuery, util.GetRaceIDFromContext(ctx))
	if err != nil {
		return
	}

	for rows.Next() {
		var gender string
		var count int
		rows.Scan(&gender, &count)
		if gender == "F" {
			female = count
		} else if gender == "M" {
			male = count
		}
	}

	return
}

func GetRaceParticipantCount(ctx context.Context, db *sqlx.DB) (int, error) {
	row := db.QueryRow(selectRaceParticipantCountQuery, util.GetRaceIDFromContext(ctx))
	if row.Err() != nil {
		return 0, row.Err()
	}

	var count int
	err := row.Scan(&count)
	return count, err
}

func GetRaceParticipantsBirthYearCount(ctx context.Context, db *sqlx.DB) ([]map[string]interface{}, error) {
	rows, err := db.Query(selectRaceParticipantBirthYearCountQuery, util.GetRaceIDFromContext(ctx))
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{} = []map[string]interface{}{}
	for rows.Next() {
		var year int
		var count int
		var gender string
		rows.Scan(&year, &gender, &count)
		results = append(results, map[string]interface{}{
			"year":   year,
			"gender": gender,
			"count":  count,
		})
	}

	return results, nil
}
