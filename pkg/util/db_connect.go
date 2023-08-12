package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
	_ "github.com/mattn/go-sqlite3"
)

func Connect(ctx context.Context, config DatabaseConfig) (*pgx.Conn, error) {
	if config.Postgres.Enabled {
		dbConn, err := connectToPostgres(ctx, config.Postgres)
		if err != nil {
			return nil, fmt.Errorf("cannot connect to db: %w", err)
		}
		return dbConn, nil
	}

	return nil, errors.New("no supported database profile enabled, please enable one (ex: postgres)")
}

func connectToPostgres(ctx context.Context, dbCfg PostgresConfig) (*pgx.Conn, error) {
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
	log.Info("Connecting to Postgres",
		"host", dbCfg.Host,
		"port", dbCfg.Port,
		"user", dbCfg.Username,
		"dbname", dbCfg.DbName,
	)
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
