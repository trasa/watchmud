package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

var watchdb *sqlx.DB

func Init() error {

	// TODO: args for db settings
	connStr := "postgres://watchmud:watchmud@localhost/watchmud?sslmode=disable"
	var err error
	watchdb, err = sqlx.Open("postgres", connStr)
	if err != nil {
		return err
	}
	if err := testConnection(); err != nil {
		return err
	}
	return nil
}

func testConnection() error {
	rows := watchdb.QueryRow("select now()")
	var now string
	if err := rows.Scan(&now); err != nil {
		return err
	}
	log.Printf("Database is live: %s", now)
	return nil
}
