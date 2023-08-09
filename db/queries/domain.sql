-- name: InsertDomain :one
INSERT INTO domains
    (name,status)
VALUES ($1,$2)
RETURNING *;

-- name: GetDomain :one
SELECT * FROM domains WHERE name = $1;

-- name: GetDomains :many
SELECT * FROM domains
ORDER BY name;

-- name: UpdateDomain :one
UPDATE domains
SET status = $1, certificate_expiry = $2, issuer = $3
WHERE name = $4
RETURNING *;

-- name: DeleteDomain :exec
DELETE FROM domains WHERE name = $1;
