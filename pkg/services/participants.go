package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
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
			birth_year,
			grade
		)
		VALUES (
			:race_id,
			:bib_number,
			:first_name,
			:last_name,
			:gender,
			:team,
			:birth_year,
			:grade
		)
	`
	selectRaceParticipantCountQuery = `
		select count(1) from participants
			where
				race_id = ?
	`

	selectRaceFinisherCountQuery = `
		select count(1) from results
			where
				race_id = ?
	`

	selectEventFinisherCountQuery = `
	select count(1) from results
		where
			race_id = ?
			and
			event_id = ?
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

	selectRaceParticipantByBibNumber = `
		select * from participants
			where race_id = ? AND bib_number = ?
	`
)

type (
	ParticipantRow struct {
		ID        int           `db:"rowid"`
		RaceID    string        `db:"race_id"`
		BibNumber string        `db:"bib_number"`
		FirstName string        `db:"first_name"`
		LastName  string        `db:"last_name"`
		BirthYear int64         `db:"birth_year"`
		Grade     sql.NullInt64 `db:"grade"`
		Gender    string        `db:"gender"`
		Team      string        `db:"team"`
	}

	ParticipantObject struct {
		ID        int    `json:"id"`
		BibNumber string `json:"bibNumber"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		BirthYear int64  `json:"birthYear,omitempty"`
		Grade     *int64 `json:"grade,omitempty"`
		Gender    string `json:"gender"`
		Team      string `json:"team"`
	}
)

func AddParticipant(ctx context.Context, db db.DB, participant ParticipantRow) error {
	_, err := db.NamedExecContext(ctx, addParticipantQuery, participant)
	return err
}

const listParticipantsQuery = `
select rowid, * from participants
	where
		race_id = ?
`

func GetParticipantByBibNumber(ctx context.Context, db db.DB, raceID, bibNumber string) (ParticipantRow, error) {
	var participant []ParticipantRow
	err := db.SelectContext(ctx, &participant, selectRaceParticipantByBibNumber, raceID, bibNumber)
	if err != nil {
		return ParticipantRow{}, err
	}

	if len(participant) != 1 {
		return ParticipantRow{}, ErrorRowNotFound
	}

	return participant[0], nil
}

func ListParticipants(ctx context.Context, db db.DB, raceID string, limit, offset int, filters map[string]string) ([]ParticipantObject, error) {
	var query string = listParticipantsQuery

	var values []interface{} = []interface{}{}
	values = append(values, raceID)
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

	query = fmt.Sprintf("%s order by bib_number*1 asc limit ? offset ?", query)
	values = append(values, limit)
	values = append(values, offset)

	log.Debug().Str("query", query).Interface("filters", filters).Send()

	var dbResults []ParticipantRow
	err := db.SelectContext(ctx, &dbResults, query, values...)
	if err != nil {
		return nil, err
	}

	results := []ParticipantObject{}
	for _, participant := range dbResults {
		pObj := ParticipantObject{
			ID:        participant.ID,
			BibNumber: participant.BibNumber,
			FirstName: participant.FirstName,
			LastName:  participant.LastName,
			Gender:    participant.Gender,
			Team:      participant.Team,
		}

		if participant.Grade.Valid {
			pObj.Grade = &participant.Grade.Int64
		}

		pObj.BirthYear = participant.BirthYear

		results = append(results, pObj)
	}

	return results, nil
}

func GetEventFinisherCount(ctx context.Context, db db.DB, raceID, eventID string) (int, error) {
	rows, err := db.QueryContext(ctx, selectEventFinisherCountQuery, raceID, eventID)
	if err != nil {
		return 0, err
	}

	for rows.Next() {
		var count int
		rows.Scan(&count)

		return count, nil
	}
	return 0, nil
}

func GetRaceFinisherCount(ctx context.Context, db db.DB, raceID string) (int, error) {
	rows, err := db.QueryContext(ctx, selectRaceFinisherCountQuery, raceID)
	if err != nil {
		return 0, err
	}

	for rows.Next() {
		var count int
		rows.Scan(&count)

		return count, err
	}
	return 0, err
}

func GetRaceGenderCount(ctx context.Context, db db.DB, raceID string) (female int, male int, err error) {
	rows, err := db.QueryContext(ctx, selectRaceParticipantGenderCountQuery, raceID)
	if err != nil {
		return
	}

	for rows.Next() {
		var gender string
		var count int
		rows.Scan(&gender, &count)
		if gender == "Female" {
			female = count
		} else if gender == "Male" {
			male = count
		}
	}

	return
}

func GetRaceParticipantCount(ctx context.Context, db db.DB, raceID string) (int, error) {
	row := db.QueryRowContext(ctx, selectRaceParticipantCountQuery, raceID)
	if row.Err() != nil {
		return 0, row.Err()
	}

	var count int
	err := row.Scan(&count)
	return count, err
}

func GetRaceParticipantsBirthYearCount(ctx context.Context, db db.DB, raceID string) ([]map[string]interface{}, error) {
	rows, err := db.QueryContext(ctx, selectRaceParticipantBirthYearCountQuery, raceID)
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

func GetRaceLastBibNumber(ctx context.Context, db db.DB) (int, error) {

	return 0, nil
}

func UpdateParticipant(ctx context.Context, db db.DB, raceID, participantID string, participant ParticipantRow) error {
	query := "UPDATE participants SET first_name = ?, last_name = ?, birth_year = ?, team = ?, gender = ?, bib_number = ? WHERE race_id = ? AND rowid = ?"
	_, err := db.ExecContext(ctx, query, participant.FirstName, participant.LastName, participant.BirthYear, participant.Team, participant.Gender, participant.BibNumber, raceID, participantID)
	return err
}
