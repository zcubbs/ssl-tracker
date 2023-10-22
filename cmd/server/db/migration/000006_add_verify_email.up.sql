CREATE  TABLE "verify_emails" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "user_id" UUID NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "secret_code" VARCHAR(32) NOT NULL,
    "is_used" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMPTZ  NOT NULL DEFAULT current_timestamp,
    "expired_at" TIMESTAMPTZ  NOT NULL DEFAULT current_timestamp + INTERVAL '15 minutes'
);

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
