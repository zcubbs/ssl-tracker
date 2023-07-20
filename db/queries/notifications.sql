-- name: InsertNotification :one
INSERT INTO notifications
  (message, send_to, channel)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetNotifications :many
SELECT * FROM notifications;

-- name: GetNotification :one
SELECT * FROM notifications WHERE id = $1;

-- name: DeleteNotification :exec
DELETE FROM notifications WHERE id = $1;

-- name: UpdateNotification :one
UPDATE notifications SET message = $1, send_to = $2, channel = $3
WHERE id = $4
RETURNING *;

-- name: GetNotificationsByChannel :many
SELECT * FROM notifications WHERE channel = $1;

