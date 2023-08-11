CREATE TABLE IF NOT EXISTS domains
(
  name TEXT PRIMARY KEY,
  certificate_expiry TIMESTAMP,
  status TEXT,
  issuer TEXT,
  owner VARCHAR(255) NOT NULL REFERENCES users(username) ON DELETE CASCADE,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  CONSTRAINT domains_name_owner UNIQUE (name, owner)
);

CREATE UNIQUE INDEX owner_idx ON "domains" ("name", "owner");
