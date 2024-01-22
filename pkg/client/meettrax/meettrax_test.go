package meettrax

import (
	"testing"
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
			got, err := getRawEventData(tt.args.meetID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRawEventData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				t.Errorf("getRawEventData() = %v, want %v", got, tt.want)
			}
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
