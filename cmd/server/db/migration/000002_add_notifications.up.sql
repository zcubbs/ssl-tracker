CREATE TABLE IF NOT EXISTS notifications
(
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "subject" TEXT NOT NULL,
  "message" TEXT NOT NULL,
  "send_to" TEXT NOT NULL,
  "channel" TEXT NOT NULL,
  "created_at"  TIMESTAMPTZ  NOT NULL DEFAULT current_timestamp
);

CREATE INDEX idx_notifications_channel ON "notifications" ("channel");
