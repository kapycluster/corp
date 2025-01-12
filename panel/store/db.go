package store

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func New(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Setup(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS control_planes (
			id TEXT PRIMARY KEY,
			name TEXT,
			user_id TEXT,
			region TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS invites (
			id TEXT PRIMARY KEY,
			used INTEGER DEFAULT 0
		);
	`
	if _, err := db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("creating table: %w", err)
	}
	return nil
}
