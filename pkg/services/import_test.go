package services

import (
	"testing"
)

func Test_convertToMilliseconds(t *testing.T) {
	type args struct {
		timing string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "test",
			args: args{
				timing: "1:03.20",
			},
			want: 63200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertToMilliseconds(tt.args.timing); got != tt.want {
				t.Errorf("convertToMilliseconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isRaceHeatResultsTab(t *testing.T) {
	type args struct {
		tab string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test",
			args: args{
				tab: "1600m A",
			},
			want: true,
		},
		{
			name: "test",
			args: args{
				tab: "400m A",
			},
			want: true,
		},
		{
			name: "test",
			args: args{
				tab: "50m",
			},
			want: true,
		},
		{
			name: "test",
			args: args{
				tab: "Race 2008-2009",
			},
			want: true,
		},
		{
			name: "test",
			args: args{
				tab: "1K Race 2008-2009",
			},
			want: true,
		},
		{
			name: "test",
			args: args{
				tab: "1K Race 2010-2011",
			},
			want: true,
		},
		{
			name: "test",
			args: args{
				tab: "1K Race 2010-2014",
			},
			want: true,
		},
		{
			name: "test",
			args: args{
				tab: "Schedule",
			},
			want: false,
		},
		{
			name: "test",
			args: args{
				tab: "Roster",
			},
			want: false,
		},
		{
			name: "test",
			args: args{
				tab: "Overall Results",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRaceHeatResultsTab(tt.args.tab); got != tt.want {
				t.Errorf("isRaceHeatResultsTab() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getHeatDistanceMeters(t *testing.T) {
	type args struct {
		description string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "test",
			args: args{
				description: "1K Race 2010-2011",
			},
			want: 1000,
		},
		{
			name: "test",
			args: args{
				description: "Kiwanis Series Sept 7, 2021",
			},
			want: 0,
		},
		{
			name: "test",
			args: args{
				description: "2022 Sept 6 Kiwanis Series 1K",
			},
			want: 1000,
		},
		{
			name: "test",
			args: args{
				description: "1600m A",
			},
			want: 1600,
		},
		{
			name: "test",
			args: args{
				description: "50m",
			},
			want: 50,
		},
		{
			name: "test",
			args: args{
				description: "400m A",
			},
			want: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHeatDistanceMeters(tt.args.description); got != tt.want {
				t.Errorf("getHeatDistanceMeters() = %v, want %v", got, tt.want)
			}
		})
	}
}
