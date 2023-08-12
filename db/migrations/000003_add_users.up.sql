CREATE TABLE IF NOT EXISTS "users"
(
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "username" VARCHAR(255) UNIQUE NOT NULL,
  "hashed_password" VARCHAR(255) NOT NULL,
  "full_name" VARCHAR(255) NOT NULL,
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "password_changed_at" TIMESTAMPTZ  NOT NULL DEFAULT current_timestamp,
  "created_at" TIMESTAMPTZ  NOT NULL DEFAULT current_timestamp
);

CREATE UNIQUE INDEX idx_users_username ON "users" ("username");
