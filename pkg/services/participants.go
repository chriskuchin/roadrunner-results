package services

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type (
	ParticipantObject struct {
		BibNumber string `json:"bibNumber"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		BirthYear int    `json:"birthYear"`
		Gender    string `json:"gender"`
		Team      string `json:"team"`
	}
)

func AddParticipant(ctx context.Context, db *sqlx.DB, participant Participant) error {
	_, err := db.NamedExec(addParticipantQuery, participant)
	return err
}

func ListParticipants(ctx context.Context, db *sqlx.DB, limit, offset int) ([]ParticipantObject, error) {
	var dbResults []Participant
	err := db.Select(&dbResults, listParticipantsQuery, raceID, limit, offset)
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

type (
	Participant struct {
		RaceID    string `db:"race_id"`
		BibNumber string `db:"bib_number"`
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		BirthYear int    `db:"birth_year"`
		Gender    string `db:"gender"`
		Team      string `db:"team"`
	}
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

	listParticipantsQuery string = `
		select * from participants
			where
				race_id = ?
			limit ? offset ?
	`
)

const (
	addEventParticipationQuery string = `
		insert into event_participation(
			race_id,
			event_id,
			bib_number
		)
		values (?,?,?)
	`
)

func RecordEventParticipation(ctx context.Context, db *sqlx.DB, raceID, eventID, bibNumber string) error {
	_, err := db.Exec(addEventParticipationQuery, raceID, eventID, bibNumber)
	return err
}
