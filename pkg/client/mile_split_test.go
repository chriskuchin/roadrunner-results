package client

import (
	"testing"
)

func Test_getEventDistanceFromHeader(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "800_meters",
			args: args{
				header: "Girls 800 meter Varsity Timed Finals",
			},
			want: 800,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEventDistanceFromHeader(tt.args.header); got != tt.want {
				t.Errorf("getEventDistanceFromHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Girls 1600 meter Varsity Timed Finals
// Boys 1600 meter Varsity Timed Finals
// Girls 400 meter Varsity Timed Finals
// Boys 400 meter Varsity Timed Finals
// Girls 800 meter Varsity Timed Finals
// Boys 800 meter Varsity Timed Finals
// Girls 3200 meter Varsity Timed Finals
// Boys 3200 meter Varsity Timed Finals

func Test_getDistanceInMeters(t *testing.T) {
	type args struct {
		distance string
		unit     string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "800 meters",
			args: args{
				distance: "800",
				unit:     "meters",
			},
			want: 800,
		},
		{
			name: "1600 meters",
			args: args{
				distance: "1600",
				unit:     "meters",
			},
			want: 1600,
		},
		{
			name: "1 kilometer",
			args: args{
				distance: "1",
				unit:     "kilometer",
			},
			want: 1000,
		},
		{
			name: "1 mile",
			args: args{
				distance: "1",
				unit:     "mile",
			},
			want: 1609.344, // 1609.334
		},
		{
			name: "5k",
			args: args{
				distance: "5",
				unit:     "kilometers",
			},
			want: 5000, // 1609.334
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDistanceInMeters(tt.args.distance, tt.args.unit); got != tt.want {
				t.Errorf("getDistanceInMeters() = %v, want %v", got, tt.want)
			}
		})
	}
}
