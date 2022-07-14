package db

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	testDBUrl = "../../results.db"
)

func setup() {
	// run migration

	// u, err := getDatabaseURL(c)
	// if err != nil {
	// 	return err
	// }
	// dbUrl, _ := url.Parse(fmt.Sprintf("sqlite:%s", testDBUrl))
	// db := dbmate.New(dbUrl)
	// db.MigrationsDir = "../../db/migrations"

	// err := db.CreateAndMigrate()
	// if err != nil {
	// 	log.Error().Err(err).Send()
	// }
}

func teardown() {
	// delete database

}

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
				db: sqlx.MustOpen("sqlite3", testDBUrl),
			},
			args: args{
				name: "2022 Spring Invitational",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RaceDAO{
				db:    tt.fields.db,
				limit: 10,
			}
			setup()
			id, _ := r.CreateRace(context.TODO(), tt.args.name)
			r.GetRaceByID(context.TODO(), id)
			r.ListRaces(context.TODO())
			r.DeleteRace(context.TODO(), id)
			teardown()
		})
	}
}
