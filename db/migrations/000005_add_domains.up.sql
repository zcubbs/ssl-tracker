CREATE TABLE IF NOT EXISTS domains
(
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" TEXT NOT NULL,
  "certificate_expiry" TIMESTAMP,
  "status" TEXT,
  "issuer" TEXT,
  "owner" UUID NOT NULL,
  "created_at" TIMESTAMPTZ  NOT NULL DEFAULT current_timestamp
);

ALTER TABLE "domains" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");
ALTER TABLE "domains" ADD CONSTRAINT "domains_name_owner" UNIQUE ("name", "owner");
CREATE INDEX idx_domains_owner ON "domains" ("owner");
