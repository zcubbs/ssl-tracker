CREATE TABLE "users" (
                       "username" varchar PRIMARY KEY,
                       "hashed_password" varchar NOT NULL,
                       "full_name" varchar NOT NULL,
                       "email" varchar UNIQUE NOT NULL,
                       "password_changed_at" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),
                       "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "domains" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

-- CREATE UNIQUE INDEX ON "domains" ("owner", "name");
ALTER TABLE "domains" ADD CONSTRAINT "owner_name_key" UNIQUE ("owner", "name");
