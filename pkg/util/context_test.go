package util

import (
	"context"
	"testing"
)

func TestGetRaceIDFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid",
			args: args{
				ctx: context.WithValue(context.Background(), RaceID, "test_id"),
			},
			want: "test_id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRaceIDFromContext(tt.args.ctx); got != tt.want {
				t.Errorf("GetRaceIDFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
