package util

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zcubbs/tlz/cmd/server/config"
)

func Connect(ctx context.Context, config config.DatabaseConfig) (*pgxpool.Pool, error) {
	if config.Postgres.Enabled {
		dbConn, err := connectToPostgres(ctx, config.Postgres)
		if err != nil {
			return nil, fmt.Errorf("cannot connect to db: %w", err)
		}
		return dbConn, nil
	}

	return nil, errors.New("no supported database profile enabled, please enable one (ex: postgres)")
}

func ConnectPostgresStdLib(dbCfg config.PostgresConfig) (*sql.DB, error) {
	dsn := getPostgresConnectionString(dbCfg)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToPostgres(ctx context.Context, dbCfg config.PostgresConfig) (*pgxpool.Pool, error) {
	dsn := getPostgresConnectionString(dbCfg)
	log.Info("connecting to Postgres",
		"host", dbCfg.Host,
		"port", dbCfg.Port,
		"user", dbCfg.Username,
		"dbname", dbCfg.DbName,
	)
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = dbCfg.MaxConns
	poolConfig.MinConns = dbCfg.MinConns
	conn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
