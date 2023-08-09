CREATE TABLE IF NOT EXISTS domains
(
    name TEXT PRIMARY KEY,
    certificate_expiry TIMESTAMP,
    status TEXT,
    issuer TEXT,
    "owner" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "domains" ("owner");
