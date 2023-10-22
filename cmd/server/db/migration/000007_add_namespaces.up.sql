CREATE TABLE IF NOT EXISTS namespaces
(
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" TEXT NOT NULL,
  "enabled" BOOLEAN NOT NULL DEFAULT TRUE,
  "billing_enabled" BOOLEAN NOT NULL DEFAULT FALSE,
  "billing_period" TEXT NOT NULL DEFAULT 'monthly',
  "billing_price" INTEGER NOT NULL DEFAULT 0,
  "billing_currency" TEXT NOT NULL DEFAULT 'usd',
  "billing_last_charge" TIMESTAMPTZ,
  "billing_next_charge" TIMESTAMPTZ,
  "owner" UUID NOT NULL,
  "created_at" TIMESTAMPTZ  NOT NULL DEFAULT current_timestamp
);

ALTER TABLE "domains" ADD COLUMN "namespace" UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000';

ALTER TABLE "namespaces" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");
ALTER TABLE "namespaces" ADD CONSTRAINT "namespaces_name_owner" UNIQUE ("name", "owner");
ALTER TABLE "domains" ADD FOREIGN KEY ("namespace") REFERENCES "namespaces" ("id");
ALTER TABLE "domains" ADD CONSTRAINT "domains_name_namespace" UNIQUE ("name", "namespace");
CREATE INDEX idx_namespaces_owner ON "namespaces" ("owner");
