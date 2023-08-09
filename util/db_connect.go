package util

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

func Connect(config DatabaseConfig) (*sql.DB, error) {
	if config.Sqlite.Enabled && config.Postgres.Enabled {
		return nil, errors.New("cannot enable both sqlite and postgres")
	}

	if !config.Sqlite.Enabled && !config.Postgres.Enabled {
		return nil, errors.New("no database profile enabled, please enable either sqlite or postgres")
	}

	if config.Sqlite.Enabled {
		dbConn, err := connectToSqlite(config.Sqlite)
		if err != nil {
			return nil, fmt.Errorf("cannot connect to db: %w", err)
		}
		return dbConn, nil
	}

	if config.Postgres.Enabled {
		dbConn, err := connectToPostgres(config.Postgres)
		if err != nil {
			return nil, fmt.Errorf("cannot connect to db: %w", err)
		}
		return dbConn, nil
	}

	return nil, errors.New("no database profile enabled, please enable either sqlite or postgres")
}

func connectToPostgres(dbCfg PostgresConfig) (*sql.DB, error) {
	var sslMode string
	if dbCfg.SslMode {
		sslMode = "enable"
	} else {
		sslMode = "disable"
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.DbName,
		sslMode,
	)
	log.Infof("Connecting to Postgres [host=%s, port=%d, dbname=%s]",
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.DbName,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToSqlite(dbCfg SqliteConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbCfg.DbName+".sqlite")
	if err != nil {
		return nil, err
	}

	return db, nil
}
