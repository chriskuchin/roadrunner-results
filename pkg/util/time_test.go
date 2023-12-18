package util

import "testing"

func Test_convertFormatTimeToMilliseconds(t *testing.T) {
	type args struct {
		formattedTime string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "2_minutes",
			args: args{
				formattedTime: "2:23.34",
			},
			want: 143340,
		},
		{
			name: "11_minutes",
			args: args{
				formattedTime: "11:05.46",
			},
			want: 665460,
		},
		{
			name: "1_minute",
			args: args{
				formattedTime: "1:00.86",
			},
			want: 60860,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertFormatTimeToMilliseconds(tt.args.formattedTime); got != tt.want {
				t.Errorf("convertFormatTimeToMilliseconds() = %v, want %v", got, tt.want)
			}
		})
	}
}
