-- name: InsertNamespace :one
INSERT INTO namespaces
    (name,owner)
VALUES ($1,$2)
RETURNING *;

-- name: GetNamespace :one
SELECT * FROM namespaces
WHERE name = $1 AND owner = $2;

-- name: GetAllNamespaces :many
SELECT * FROM namespaces
WHERE owner = $1
ORDER BY name;
