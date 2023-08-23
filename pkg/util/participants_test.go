package util

import "testing"

func TestSplitName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "simple_name",
			args: args{
				name: "John Smith",
			},
			want:  "John",
			want1: "Smith",
		},
		{
			name: "first_name_only",
			args: args{
				name: "John",
			},
			want:  "John",
			want1: "",
		},
		{
			name: "multi_name",
			args: args{
				name: "James de la hoya",
			},
			want:  "James",
			want1: "de la hoya",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SplitName(tt.args.name)
			if got != tt.want {
				t.Errorf("SplitName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SplitName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
