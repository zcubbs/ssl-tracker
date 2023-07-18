package db

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zcubbs/tlz/util"
)

var Store Storer

type Storer interface {
	Querier
	GetConn() *sql.DB
}

type SQLiteDB struct {
	*Queries
	conn *sql.DB
}

type PostgresDB struct {
	*Queries
	conn *sql.DB
}

func Connect(config util.DatabaseConfig) {
	if config.Sqlite.Enabled && config.Postgres.Enabled {
		log.Fatal("cannot enable both sqlite and postgres")
	}

	if !config.Sqlite.Enabled && !config.Postgres.Enabled {
		log.Fatal("no database profile enabled, please enable either sqlite or postgres")
	}

	if config.Sqlite.Enabled {
		conn, err := connectToSqlite(config.Sqlite)
		if err != nil {
			log.Fatal("cannot connect to db:", err)
		}
		Store = &SQLiteDB{
			conn:    conn,
			Queries: New(conn),
		}
		log.Info("Connected to SQLite")
		return
	}

	if config.Postgres.Enabled {
		conn, err := connectToPostgres(config.Postgres)
		if err != nil {
			log.Fatal("cannot connect to db:", err)
		}
		Store = &PostgresDB{
			conn:    conn,
			Queries: New(conn),
		}
		log.Info("Connected to Postgres")
		return
	}
}

func (db *SQLiteDB) GetConn() *sql.DB {
	return db.conn
}

func (db *PostgresDB) GetConn() *sql.DB {
	return db.conn
}

func connectToPostgres(dbCfg util.PostgresConfig) (*sql.DB, error) {
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

func connectToSqlite(dbCfg util.SqliteConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbCfg.DbName+".sqlite")
	if err != nil {
		return nil, err
	}

	return db, nil
}
