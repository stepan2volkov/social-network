package userrepo

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
)

func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error to open db connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error to ping db connection: %w", err)
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("error to set mysql dialect for migrations: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("error to migrate: %w", err)
	}

	return nil
}
