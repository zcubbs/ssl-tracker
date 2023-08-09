package db

import (
	"database/sql"
)

type Store struct {
	*Queries
	db *sql.DB
}

func (s Store) GetConn() *sql.DB {
	return s.db
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}
