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
			name: "26-minute-change",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateMilisecondsFromTiming(tt.args.min, tt.args.sec, tt.args.hund); got != tt.want {
				t.Errorf("calculateMilisecondsFromTiming() = %v, want %v", got, tt.want)
			}
		})
	}
}
