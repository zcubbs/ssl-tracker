CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message TEXT NOT NULL,
    send_to TEXT NOT NULL,
    channel TEXT NOT NULL,
    created_at  TIMESTAMP  NOT NULL DEFAULT current_timestamp
)
