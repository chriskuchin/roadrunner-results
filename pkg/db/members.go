package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type (
	Member struct {
		RaceID    string `db:"race_id"`
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		BirthYear int    `db:"birth_year"`
		Gender    string `db:"gender"`
	}
)

type MemberDao struct {
	db *sqlx.DB
}

func NewMemberDAO(db *sqlx.DB) *MemberDao {
	return &MemberDao{
		db: db,
	}
}

func (m *MemberDao) UpsertMember(ctx context.Context, firstName, lastName, gender string, birthYear int) error {
	_, err := m.db.Exec("insert into members (first_name, last_name, birth_year, gender) VALUES(?, ?, ?, ?)", firstName, lastName, birthYear, gender)
	return err
}
