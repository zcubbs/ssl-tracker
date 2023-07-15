package main

import (
	"database/sql"
)

type Domain struct {
	Name              string       `json:"name"`
	CertificateExpiry sql.NullTime `json:"certificate_expiry"`
}
