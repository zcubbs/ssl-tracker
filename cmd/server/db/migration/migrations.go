package migration

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	mig "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/mattes/migrate/source/file"
	"github.com/zcubbs/ssl-tracker/cmd/server/config"
	dbConnect "github.com/zcubbs/ssl-tracker/cmd/server/db/connect"
	"github.com/zcubbs/ssl-tracker/cmd/server/logger"
	"net/http"
)

var (
	log = logger.L()
)

//go:embed *.sql
var migrations embed.FS

func Run(dbCfg config.DatabaseConfig) error {
	if dbCfg.Postgres.Enabled {
		conn, err := dbConnect.ConnectPostgresStdLib(dbCfg.Postgres)
		if err != nil {
			return err
		}
		defer func(conn *sql.DB) {
			err := conn.Close()
			if err != nil {
				log.Fatal("failed to close database connection", "error", err)
			}
		}(conn)

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

	m, err := mig.NewWithInstance(
		"FS",
		source,
		dbname,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, mig.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Info("✔️ applied database migrations")
	return nil
}
