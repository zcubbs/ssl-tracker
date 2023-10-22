CREATE TABLE IF NOT EXISTS domains
(
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" TEXT NOT NULL,
  "certificate_expiry" TIMESTAMP,
  "status" TEXT,
  "issuer" TEXT,
  "created_at" TIMESTAMPTZ  NOT NULL DEFAULT current_timestamp
);
