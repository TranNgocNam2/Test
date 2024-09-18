package db

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func ConnectDB(dsn string) (db *sqlx.DB, err error) {
	db, err = sqlx.Connect("pgx", dsn)

	if err != nil {
		return nil, err
	}
	return db, nil
}
