package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zcubbs/logwrapper/logger"
)

func DbConnect(ctx context.Context, config DatabaseConfig, logger logger.Logger) (*pgxpool.Pool, error) {
	if config.Postgres.Enabled {
		dbConn, err := connectToPostgres(ctx, config.Postgres, logger)
		if err != nil {
			return nil, fmt.Errorf("cannot connect to db: %w", err)
		}
		return dbConn, nil
	}

	return nil, errors.New("no supported database profile enabled, please enable one (ex: postgres)")
}

func connectToPostgres(ctx context.Context, dbCfg PostgresConfig, logger logger.Logger) (*pgxpool.Pool, error) {
	dsn := getPostgresConnectionString(dbCfg)
	logger.Info("connecting to Postgres",
		"host", dbCfg.Host,
		"port", dbCfg.Port,
		"user", dbCfg.Username,
		"dbname", dbCfg.DbName,
	)
	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func GetDbConnectionString(config DatabaseConfig) string {
	if config.Postgres.Enabled {
		return getPostgresConnectionString(config.Postgres)
	}

	return ""
}

func getPostgresConnectionString(dbCfg PostgresConfig) string {
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

	return dsn
}
