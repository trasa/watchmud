package db

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/serverconfig"
)

var watchdb *sqlx.DB

func Init(config *serverconfig.Config) error {

	if config.DB.UseSSH {
		return errors.New("SSH not implemented yet")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DB.User,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name)
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
	log.Info().Msgf("Database is live: %s", now)
	return nil
}
