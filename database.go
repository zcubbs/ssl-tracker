package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	AddDomain(domain string) error
	GetDomains() ([]Domain, error)
	UpdateDomain(domain Domain) error
}

type SQLiteDB struct {
	conn *sql.DB
}

func NewSQLiteDB(filepath string) *SQLiteDB {
	// Open the SQLite database
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalf("Cannot open SQLite database: %v", err)
	}

	// Create the domains table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS domains 
		(name TEXT PRIMARY KEY, certificate_expiry TIMESTAMP)`)
	if err != nil {
		log.Fatalf("Cannot create domains table: %v", err)
	}

	return &SQLiteDB{conn: db}
}

func (db *SQLiteDB) AddDomain(domain string) error {
	// Assuming `db.conn` is your database connection
	_, err := db.conn.Exec("INSERT INTO domains (name) VALUES (?)", domain)
	return err
}

func (db *SQLiteDB) GetDomains() ([]Domain, error) {
	rows, err := db.conn.Query("SELECT name, certificate_expiry FROM domains")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []Domain
	for rows.Next() {
		var domain Domain
		if err := rows.Scan(&domain.Name, &domain.CertificateExpiry); err != nil {
			return nil, err
		}
		domains = append(domains, domain)
	}

	return domains, nil
}

func (db *SQLiteDB) UpdateDomain(domain Domain) error {
	_, err := db.conn.Exec("UPDATE domains SET certificate_expiry = ? WHERE name = ?", domain.CertificateExpiry, domain.Name)
	return err
}

type PostgresDB struct {
	// Postgres connection
}

func (db *PostgresDB) AddDomain(domain string) error {
	// Implement this
	return nil
}

func (db *PostgresDB) GetDomains() ([]Domain, error) {
	// Implement this
	return nil, nil
}

func (db *PostgresDB) UpdateDomain(domain Domain) error {
	// Implement this
	return nil
}
