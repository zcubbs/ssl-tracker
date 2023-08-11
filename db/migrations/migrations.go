package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/mattes/migrate/source/file"
	"github.com/zcubbs/tlz/util"
	"net/http"
)

//go:embed *.sql
var migrations embed.FS

func Migrate(conn *sql.DB, databaseType util.DatabaseType) error {
	if databaseType == util.Sqlite {
		log.Info("SQLite database enabled")
		return migrateSqlite(conn)
	}

	if databaseType == util.Postgres {
		log.Info("Postgres database enabled")
		return migratePostgres(conn)
	}

	return fmt.Errorf("unknown database type: %s", databaseType)
}

func migratePostgres(conn *sql.DB) error {
	driver, err := postgres.WithInstance(conn, &postgres.Config{
		DatabaseName:          "postgres",
		SchemaName:            "public",
		MultiStatementEnabled: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	source, err := httpfs.New(http.FS(migrations), ".")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"FS",
		source,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migration: %w", err)
	}

	log.Info("database migration completed")
	return nil
}

func migrateSqlite(conn *sql.DB) error {
	instance, err := sqlite3.WithInstance(conn, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	source, err := httpfs.New(http.FS(migrations), ".")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"FS",
		source,
		"sqlite3",
		instance,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migration: %w", err)
	}

	log.Info("database migration completed")
	return nil
}
