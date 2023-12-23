package mile_split

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/chriskuchin/roadrunner-results/pkg/client/model"
)

func Test_processRawEventData(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "full data",
			args: args{
				filename: "mile_split.txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(fmt.Sprintf("testdata/%s", tt.args.filename))
			if err != nil {
				t.Errorf("failed to open testdata file: %v", err)
			}

			data, err := io.ReadAll(f)
			if err != nil {
				t.Errorf("failed to read tesdata file: %v", err)
			}
			processRawEventData(string(data))
		})
	}
}

func Test_getResultFromRow(t *testing.T) {
	type args struct {
		row string
	}
	tests := []struct {
		name    string
		args    args
		want    model.Result
		wantErr bool
	}{
		{
			name: "first_minute_plus",
			args: args{
				row: "  1  Paul Scown                Sr Stansbury                   4:16.85",
			},
			want: model.Result{
				Place:     1,
				FirstName: "Paul",
				LastName:  "Scown",
				Year:      "Sr",
				Team:      "Stansbury",
				Time:      256850,
			},
			wantErr: false,
		},
		{
			name: "sub_minute",
			args: args{
				row: "   1  Sarah Ballard             Sr Gold Medal Athletics          56.54",
			},
			want: model.Result{
				Place:     1,
				FirstName: "Sarah",
				LastName:  "Ballard",
				Year:      "Sr",
				Team:      "Gold Medal Athletics",
				Time:      56540,
			},
			wantErr: false,
		},
		{
			name: "no_year",
			args: args{
				row: " 72  Matt Hoopes                  T-BIRDS XC                    59.80",
			},
			want: model.Result{
				Place:     72,
				FirstName: "Matt",
				LastName:  "Hoopes",
				Year:      "",
				Team:      "T-BIRDS XC",
				Time:      59800,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getResultFromRow(tt.args.row, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("getResultFromRow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getResultFromRow() = %v, want %v", got, tt.want)
			}
		})
	}
}
