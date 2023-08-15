package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/mattes/migrate/source/file"
	"github.com/zcubbs/tlz/internal/util"
	"net/http"
)

//go:embed *.sql
var migrations embed.FS

func Run(dbCfg util.DatabaseConfig) error {
	conn, err := connectToDb(dbCfg)
	if err != nil {
		return err
	}

	if dbCfg.Postgres.Enabled {
		err = migratePostgres(conn, dbCfg.Postgres.DbName)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("no database type specified")
}

func migratePostgres(conn *sql.DB, dbname string) error {
	driver, err := postgres.WithInstance(conn, &postgres.Config{
		DatabaseName:          dbname,
		SchemaName:            "public",
		MultiStatementEnabled: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	source, err := httpfs.New(http.FS(migrations), ".")
	if err != nil {
		return fmt.Errorf("failed to create migrations source: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"FS",
		source,
		dbname,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Info("✔️ applied database migrations")
	return nil
}

func connectToDb(dbCfg util.DatabaseConfig) (*sql.DB, error) {
	dsn := util.GetDbConnectionString(dbCfg)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
