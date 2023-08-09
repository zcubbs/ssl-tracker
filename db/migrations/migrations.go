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
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/util"
	"net/http"
)

//go:embed *.sql
var migrations embed.FS

func Migrate(store *db.Store, databaseType util.DatabaseType) error {
	if databaseType == util.Postgres {
		return migrateSqlite(store.GetConn())
	}

	if databaseType == util.Sqlite {
		return migratePostgres(store.GetConn())
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
