package services

import (
	"testing"
)

func Test_parseFinishingTime(t *testing.T) {
	type args struct {
		time string
	}
	tests := []struct {
		name           string
		args           args
		wantMin        int
		wantSec        int
		wantHundredths int
	}{
		{
			name: "one-minute",
			args: args{
				time: "1:00:00",
			},
			wantMin:        1,
			wantSec:        0,
			wantHundredths: 0,
		},
		{
			name: "one-second",
			args: args{
				time: "00:01:00",
			},
			wantMin:        0,
			wantSec:        1,
			wantHundredths: 0,
		},
		{
			name: "one-hundredth",
			args: args{
				time: "00:00:01",
			},
			wantMin:        0,
			wantSec:        0,
			wantHundredths: 1,
		},
		{
			name: "realistic",
			args: args{
				time: "26:24:51",
			},
			wantMin:        26,
			wantSec:        24,
			wantHundredths: 51,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotSec, gotHundredths := parseFinishingTime(tt.args.time)
			if gotMin != tt.wantMin {
				t.Errorf("parseFinishingTime() gotMin = %v, want %v", gotMin, tt.wantMin)
			}
			if gotSec != tt.wantSec {
				t.Errorf("parseFinishingTime() gotSec = %v, want %v", gotSec, tt.wantSec)
			}
			if gotHundredths != tt.wantHundredths {
				t.Errorf("parseFinishingTime() gotHundredths = %v, want %v", gotHundredths, tt.wantHundredths)
			}
		})
	}
}

func Test_calculateMilisecondsFromTiming(t *testing.T) {
	type args struct {
		min  int
		sec  int
		hund int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "one-minute",
			args: args{
				min:  1,
				sec:  0,
				hund: 0,
			},
			want: 60000,
		},
		{
			name: "one-second",
			args: args{
				min:  0,
				sec:  1,
				hund: 0,
			},
			want: 1000,
		},
		{
			name: "one-hundredths",
			args: args{
				min:  0,
				sec:  0,
				hund: 01,
			},
			want: 10,
		},
		{
			name: "realistic",
			args: args{
				min:  26,
				sec:  32,
				hund: 51,
			},
			want: 1592510,
		},
		{
			name: "realistic",
			args: args{
				min:  26,
				sec:  02,
				hund: 01,
			},
			want: 1562010,
		},
		{
			name: "realistic",
			args: args{
				min:  0,
				sec:  02,
				hund: 01,
			},
			want: 2010,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateMilisecondsFromTiming(tt.args.min, tt.args.sec, tt.args.hund); got != tt.want {
				t.Errorf("calculateMilisecondsFromTiming() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatFinishingTime(t *testing.T) {
	type args struct {
		ms int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "realistic",
			args: args{
				ms: 1592510,
			},
			want: "26:32:51",
		},
		{
			name: "leading-zeros",
			args: args{
				ms: 1562010,
			},
			want: "26:02:01",
		},
		{
			name: "zero-minutes",
			args: args{
				ms: 2010,
			},
			want: "00:02:01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatFinishingTime(tt.args.ms); got != tt.want {
				t.Errorf("formatFinishingTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
