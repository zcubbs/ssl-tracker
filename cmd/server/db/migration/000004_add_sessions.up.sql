CREATE TABLE IF NOT EXISTS "sessions"
(
  "id" UUID PRIMARY KEY,
  "user_id" UUID NOT NULL,
  "refresh_token" TEXT NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT current_timestamp
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
