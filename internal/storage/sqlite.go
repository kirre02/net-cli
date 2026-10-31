package storage

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

type SQLite struct {
	db *sql.DB
}

func NewSQLite() (*SQLite, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("get user config directory: %w", err)
	}

	appDir := filepath.Join(configDir, "net-cli")

	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, fmt.Errorf("create application directory: %w", err)
	}

	dbPath := filepath.Join(appDir, "net-cli.db")

	db, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(1)

	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("run database migrations: %w", err)
	}

	return &SQLite{db: db}, nil
}

func migrate(db *sql.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("set migration dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}
