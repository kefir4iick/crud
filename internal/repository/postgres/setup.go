package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	createTableQuery = `
		CREATE TABLE IF NOT EXISTS cars (
			id VARCHAR(36) PRIMARY KEY,
			make VARCHAR(255) NOT NULL,
			model VARCHAR(255) NOT NULL,
			year INTEGER NOT NULL,
			price INTEGER NOT NULL
		)
	`
)

func NewDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
}
