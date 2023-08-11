CREATE TABLE IF NOT EXISTS "sessions" (
                       "id" uuid PRIMARY KEY,
                       "username" varchar NOT NULL REFERENCES "users" ("username") ON DELETE CASCADE,
                       "refresh_token" varchar NOT NULL,
                       "user_agent" varchar NOT NULL,
                       "client_ip" varchar UNIQUE NOT NULL,
                       "is_blocked" boolean NOT NULL DEFAULT false,
                       "expires_at" timestamptz NOT NULL,
                       "created_at" timestamptz NOT NULL DEFAULT (now())
);
