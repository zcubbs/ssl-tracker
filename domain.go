package main

import (
	"database/sql"
)

type Status string

const (
	StatusValid    Status = "valid"
	StatusExpired  Status = "expired"
	StatusUnknown  Status = "unknown"
	StatusExpiring Status = "expiring"
	StatusPending  Status = "pending"
)

type Domain struct {
	Name              string         `json:"name"`
	Status            Status         `json:"status"`
	Issuer            sql.NullString `json:"issuer"`
	CertificateExpiry sql.NullTime   `json:"certificate_expiry"`
}

type DomainWrapper struct {
	Domain
	Until string `json:"until"`
}
