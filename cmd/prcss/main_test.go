package main

import "testing"

func Test_convertToMilliseconds(t *testing.T) {
	type args struct {
		timing string
	}
	tests := []struct {
		name string
		args args
		want int
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
