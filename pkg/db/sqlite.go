package db

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

type SQLiteDB struct {
	read  *sqlx.DB
	write *sqlx.DB
}

const (
	dbURL = "%s?_journal_mode=WAL&mode=rwc&cache=shared&_txlock=immediate&_busy_timeout=5000&_synchronous=NORMAL&_cache_size=1000000000&_foreign_keys=true"
)

func NewSQLite(dbPath string) *SQLiteDB {

	db := &SQLiteDB{}
	log.Debug().Str("db", dbPath).Send()
	write, err := sqlx.Open("sqlite3", fmt.Sprintf(dbURL, dbPath))
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	write.SetMaxOpenConns(1)

	configureSQLite(write)

	read, err := sqlx.Open("sqlite3", fmt.Sprintf(dbURL, dbPath))
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	read.SetMaxOpenConns(max(4, runtime.NumCPU()))
	configureSQLite(read)

	db.write = write
	db.read = read

	return db
}

func configureSQLite(db *sqlx.DB) (err error) {
	pragmas := []string{
		// "busy_timeout = 5000",
		// "synchronous = NORMAL",
		// "cache_size = 1000000000", // 1GB
		// "foreign_keys = true",
		"temp_store = memory",
		// "mmap_size = 3000000000",
	}

	for _, pragma := range pragmas {
		_, err = db.Exec("PRAGMA " + pragma)
		if err != nil {
			return
		}
	}

	return nil
}

func (db *SQLiteDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.write.ExecContext(ctx, query, args...)
}

func (db *SQLiteDB) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return db.write.NamedExecContext(ctx, query, arg)
}

func (db *SQLiteDB) SelectContext(ctx context.Context, destination interface{}, query string, args ...interface{}) error {
	return db.read.SelectContext(ctx, destination, query, args...)
}

func (db *SQLiteDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.read.QueryRowContext(ctx, query, args...)
}

func (db *SQLiteDB) GetContext(ctx context.Context, destination interface{}, query string, args ...interface{}) error {
	return db.read.GetContext(ctx, destination, query, args...)
}

func (db *SQLiteDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.read.QueryContext(ctx, query, args...)
}
