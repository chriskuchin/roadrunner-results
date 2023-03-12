package db

import "github.com/jmoiron/sqlx"

var dbInstance *sqlx.DB

func Configure(db *sqlx.DB) {
	dbInstance = db
}

func getDBInstance() *sqlx.DB {
	if dbInstance != nil {
		return dbInstance
	}

	return nil
}
