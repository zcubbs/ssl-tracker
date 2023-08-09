CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    subject TEXT NOT NULL,
    message TEXT NOT NULL,
    send_to TEXT NOT NULL,
    channel TEXT NOT NULL,
    created_at  TIMESTAMP  NOT NULL DEFAULT current_timestamp
);

CREATE INDEX channel_idx ON "notifications" ("channel");
