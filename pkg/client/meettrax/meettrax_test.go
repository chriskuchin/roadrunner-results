package meettrax

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/chriskuchin/roadrunner-results/pkg/client/model"
)

func Test_getRawEventData(t *testing.T) {
	type args struct {
		meetID string
	}
	tests := []struct {
		name    string
		args    args
		want    []MTTRXRawResults
		wantErr bool
	}{
		{
			name: "byu",
			args: args{
				meetID: "19",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("Comment out to generate new Data")
			got, _ := getRawEventData(tt.args.meetID)

			data, _ := json.Marshal(got)
			os.WriteFile("meettrax_raw.json", data, 0667)
		})
	}
}

func Test_getMeetIDFromURL(t *testing.T) {
	type args struct {
		eventURL string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "root_url",
			args: args{
				eventURL: "https://meettrax.com/meets/19/",
			},
			want: "19",
		},
		{
			name: "events_url",
			args: args{
				eventURL: "https://meettrax.com/meets/19/results/by-event",
			},
			want: "19",
		},
		{
			name: "event_results_url",
			args: args{
				eventURL: "https://meettrax.com/meets/19/results/by-event/741",
			},
			want: "19",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMeetIDFromURL(tt.args.eventURL); got != tt.want {
				t.Errorf("getMeetIDFromURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processRawEventData(t *testing.T) {
	tests := []struct {
		name string
		want []model.Event
	}{
		{
			name: "meettrax_byu_data",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testdata, err := os.ReadFile(fmt.Sprintf("testdata/%s.json", tt.name))
			if err != nil {
				t.Errorf("processRawEventData() failed to read testdata")
			}

			var testInput []MTTRXRawResults
			err = json.Unmarshal(testdata, &testInput)
			if err != nil {
				t.Errorf("processRawEventData() = %v", err)
			}

			processRawEventData(testInput)
			// if got := processRawEventData(testInput); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("processRawEventData() = %v, want %v", got, tt.want)
			// }
		})
	}
}
