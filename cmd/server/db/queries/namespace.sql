-- name: InsertNamespace :one
INSERT INTO namespaces
    (name,user_id)
VALUES ($1,$2)
RETURNING *;

-- name: GetNamespace :one
SELECT * FROM namespaces
WHERE name = $1 AND user_id = $2;

-- name: GetAllNamespaces :many
SELECT * FROM namespaces
WHERE user_id = $1
ORDER BY name;
