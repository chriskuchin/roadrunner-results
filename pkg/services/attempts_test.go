package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func dbSetup(dbName string) *sqlx.DB {
	dbPath := fmt.Sprintf("./testdata/%s.db", dbName)
	cmd := exec.Command("dbmate", fmt.Sprintf("--url=sqlite:%s", dbPath), "--migrations-dir=../../db/migrations", "--no-dump-schema", "up")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal().Err(err).Msgf("%s", out)
	}

	db, err := sqlx.Open("sqlite3", fmt.Sprintf("%s?mode=rwc&cache=shared", dbPath))
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return db
}

func dbCleanup(dbName string) {
	os.Remove(fmt.Sprintf("./testdata/%s.db", dbName))
}

func TestRecordAttempt(t *testing.T) {
	type args struct {
		ctx       context.Context
		raceID    string
		eventID   string
		bib       string
		attemptNo int
		result    float32
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "create_simple",
			args: args{
				ctx:       nil,
				raceID:    "test",
				eventID:   "test_event",
				bib:       "123",
				attemptNo: 1,
				result:    120.43,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		db := dbSetup(tt.name)
		t.Run(tt.name, func(t *testing.T) {
			if err := RecordAttempt(tt.args.ctx, db, tt.args.raceID, tt.args.eventID, tt.args.bib, tt.args.attemptNo, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("RecordAttempt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		dbCleanup(tt.name)
	}
}

func TestUpsertBestAttemptResults(t *testing.T) {
	type row struct {
		race_id  string
		event_id string
		bib      string
		attempt  int
		result   float32
	}
	type args struct {
		race_id  string
		event_id string
		bib      string
	}
	tests := []struct {
		rows    []row
		args    args
		insert  row
		name    string
		wantErr bool
	}{
		{
			name: "pre_existing",
			rows: []row{
				{
					race_id:  "123",
					event_id: "453",
					bib:      "444",
					attempt:  1,
					result:   120.9,
				},
				{
					race_id:  "123",
					event_id: "453",
					bib:      "444",
					attempt:  2,
					result:   120.7,
				},
				{
					race_id:  "123",
					event_id: "453",
					bib:      "444",
					attempt:  3,
					result:   121,
				},
			},
			args: args{
				race_id:  "123",
				event_id: "453",
				bib:      "444",
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		db := dbSetup(tt.name)
		t.Run(tt.name, func(t *testing.T) {
			for _, row := range tt.rows {
				RecordAttempt(context.TODO(), db, row.race_id, row.event_id, row.bib, row.attempt, row.result)
			}

			if err := UpsertBestAttemptResults(context.TODO(), db, tt.args.race_id, tt.args.event_id, tt.args.bib); (err != nil) != tt.wantErr {
				t.Errorf("UpsertBestAttemptResults() error = %v, wantErr %v", err, tt.wantErr)
			}

			ListAttempts(context.TODO(), db, tt.args.race_id, tt.args.event_id, tt.args.bib)
		})
		dbCleanup(tt.name)
	}
}
