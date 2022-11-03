package db

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestEventDao_AddEvent(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EventDao{
				db: tt.fields.db,
			}
			e.AddEvent(tt.args.ctx)
		})
	}
}
