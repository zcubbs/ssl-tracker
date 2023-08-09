CREATE TABLE IF NOT EXISTS domains
(
    name TEXT PRIMARY KEY,
    certificate_expiry TIMESTAMP,
    status TEXT,
    issuer TEXT,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX owner_idx ON "domains" ("owner");
