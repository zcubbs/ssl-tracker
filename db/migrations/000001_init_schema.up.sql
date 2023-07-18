CREATE TABLE IF NOT EXISTS domains
(
    name TEXT PRIMARY KEY,
    certificate_expiry TIMESTAMP,
    status TEXT,
    issuer TEXT
)
