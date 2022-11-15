package services

import (
	"context"
	"testing"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	testDBUrl = "../../results.db"
	raceDAO   = db.NewRaceDAO(sqlx.MustOpen("sqlite3", testDBUrl))
)

func setup() context.Context {
	NewRaceService(raceDAO)
	raceID, _ := GetRaceServiceInstance().CreateRace(context.TODO(), "test race 2022")
	ctx := context.WithValue(context.Background(), util.RaceID, raceID)

	return ctx
}

func TestParticipantsService_ProcessParticipantCSV(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		ps   *ParticipantsService
		args args
	}{
		{
			name: "basic",
			ps:   &ParticipantsService{},
			args: args{
				filePath: "testdata/participants-test.csv",
			},
		},
	}
	for _, tt := range tests {
		ctx := setup()
		t.Run(tt.name, func(t *testing.T) {
			ps := &ParticipantsService{}
			ps.ProcessParticipantCSV(ctx, tt.args.filePath)
		})
	}
}
