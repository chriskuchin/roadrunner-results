package util

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
			name: "3000_meters",
			args: args{
				header: "3,000 Meters Middle School 6th - FRI SEPT 22 6:00p",
			},
			want: 3000,
		},
		{
			name: "800_meters",
			args: args{
				header: "Girls 800 meter Varsity Timed Finals",
			},
			want: 800,
		},
		{
			name: "1600_meter_girls",
			args: args{
				header: "Girls 1600 meter Varsity Timed Finals",
			},
			want: 1600,
		},
		{
			name: "1600_meter_boys",
			args: args{
				header: "Boys 1600 meter Varsity Timed Finals",
			},
			want: 1600,
		},
		{
			name: "3200_meter_boys",
			args: args{
				header: "Boys 3200 meter Varsity Timed Finals",
			},
			want: 3200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEventDistanceFromHeader(tt.args.header); got != tt.want {
				t.Errorf("getEventDistanceFromHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			want: 5000,
		},
		{
			name: "1 m",
			args: args{
				distance: "1",
				unit:     "m",
			},
			want: 1609.344,
		},
		{
			name: "100 m",
			args: args{
				distance: "100",
				unit:     "m",
			},
			want: 100,
		},
		{
			name: "marathon",
			args: args{
				distance: "26.2",
				unit:     "m",
			},
			want: 42164.8128,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDistanceInMeters(tt.args.distance, tt.args.unit); got != tt.want {
				t.Errorf("getDistanceInMeters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGenderFromHeader(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "800_meters",
			args: args{
				header: "Girls 800 meter Varsity Timed Finals",
			},
			want: "Female",
		},
		{
			name: "1600_meter_girls",
			args: args{
				header: "Girls 1600 meter Varsity Timed Finals",
			},
			want: "Female",
		},
		{
			name: "1600_meter_boys",
			args: args{
				header: "Boys 1600 meter Varsity Timed Finals",
			},
			want: "Male",
		},
		{
			name: "3200_meter_boys",
			args: args{
				header: "Boys 3200 meter Varsity Timed Finals",
			},
			want: "Male",
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetGenderFromHeader(tt.args.header); got != tt.want {
				t.Errorf("GetGenderFromHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
