package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func dbSetup(dbName string) db.DB {
	dbPath := fmt.Sprintf("./testdata/%s.db", dbName)
	cmd := exec.Command("dbmate", fmt.Sprintf("--url=sqlite:%s", dbPath), "--migrations-dir=../../db/migrations", "--no-dump-schema", "up")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal().Err(err).Msgf("%s", out)
	}

	return db.NewSQLite(dbPath)
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

func TestListAttempts(t *testing.T) {
	type row struct {
		race_id  string
		event_id string
		bib      string
		attempt  int
		result   float32
	}
	type args struct {
		raceID  string
		eventID string
		bib     string
	}
	tests := []struct {
		name    string
		args    args
		rows    []row
		want    []AttemptsRow
		wantErr bool
	}{
		{
			name: "list",
			args: args{
				raceID:  "123",
				eventID: "453",
				bib:     "444",
			},
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
			want: []AttemptsRow{
				{
					RaceID:  "123",
					EventID: "453",
					Bib:     "444",
					Attempt: 3,
					Result:  121,
				},
				{
					RaceID:  "123",
					EventID: "453",
					Bib:     "444",
					Attempt: 1,
					Result:  120.9,
				},
				{
					RaceID:  "123",
					EventID: "453",
					Bib:     "444",
					Attempt: 2,
					Result:  120.7,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := dbSetup(tt.name)
			for _, row := range tt.rows {
				RecordAttempt(context.TODO(), db, row.race_id, row.event_id, row.bib, row.attempt, row.result)
			}
			got, err := ListAttempts(context.TODO(), db, tt.args.raceID, tt.args.eventID, tt.args.bib)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAttempts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for idx, attempt := range got {
				if !reflect.DeepEqual(attempt, tt.want[idx]) {
					t.Errorf("ListAttempts() = %v, want %v", attempt, tt.want[idx])
				}
			}
		})
	}
}
