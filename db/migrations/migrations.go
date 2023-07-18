package migrations

import (
	"database/sql"
	"embed"
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

func Migrate(dbCfg util.DatabaseConfig) {
	if dbCfg.Sqlite.Enabled {
		migrateSqlite(db.Store.GetConn())
		return
	}

	if dbCfg.Postgres.Enabled {
		migratePostgres(db.Store.GetConn())
		return
	}
}

func migratePostgres(conn *sql.DB) {
	driver, err := postgres.WithInstance(conn, &postgres.Config{
		DatabaseName:          "postgres",
		SchemaName:            "public",
		MultiStatementEnabled: true,
	})
	if err != nil {
		log.Fatal(err, "failed to create postgres driver")
	}

	source, err := httpfs.New(http.FS(migrations), ".")
	if err != nil {
		log.Fatal(err, "failed to create migration source")
	}

	m, err := migrate.NewWithInstance(
		"FS",
		source,
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err, "failed to create migration")
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err, "failed to apply migration")
	}

	log.Info("database migration completed")
}

func migrateSqlite(conn *sql.DB) {
	instance, err := sqlite3.WithInstance(conn, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err, "failed to create postgres driver")
	}

	source, err := httpfs.New(http.FS(migrations), ".")
	if err != nil {
		log.Fatal(err, "failed to create migration source")
	}

	m, err := migrate.NewWithInstance(
		"FS",
		source,
		"sqlite3",
		instance,
	)
	if err != nil {
		log.Fatal(err, "failed to create migration")
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err, "failed to apply migration")
	}

	log.Info("database migration completed")
}
