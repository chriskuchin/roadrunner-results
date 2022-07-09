package db

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func TestRaceDAO_CreateRace(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test",
			fields: fields{
				db: sqlx.MustOpen("sqlite3", "../../results.db"),
			},
			args: args{
				name: "2022 Spring Invitational",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RaceDAO{
				db: tt.fields.db,
			}
			r.CreateRace(tt.args.name)
		})
	}
}
