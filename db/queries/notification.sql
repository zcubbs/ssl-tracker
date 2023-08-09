-- name: InsertNotification :one
INSERT INTO notifications
  (subject, message, send_to, channel)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetNotifications :many
SELECT * FROM notifications;

-- name: GetNotification :one
SELECT * FROM notifications WHERE id = $1;

-- name: DeleteNotification :exec
DELETE FROM notifications WHERE id = $1;

-- name: UpdateNotification :one
UPDATE notifications SET subject = $1, message = $2, send_to = $3, channel = $4
WHERE id = $5
RETURNING *;

-- name: GetNotificationsByChannel :many
SELECT * FROM notifications WHERE channel = $1;

