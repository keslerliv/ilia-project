package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/keslerliv/user/config"
)

func OpenConnection() (*sql.DB, error) {
	sc := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Env.DBHost,
		config.Env.DBPort,
		config.Env.DBUser,
		config.Env.DBPassword,
		config.Env.DBName,
		config.Env.DBSSLMode,
	)

	conn, err := sql.Open("postgres", sc)
	if err != nil {
		panic(err)
	}

	err = conn.Ping()

	return conn, err
}

func MakeMigration(conn *sql.DB) error {
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	return m.Up()
}
