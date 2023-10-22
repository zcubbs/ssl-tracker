-- name: InsertDomain :one
INSERT INTO domains
    (name,status,namespace)
VALUES ($1,$2,$3)
RETURNING *;

-- name: GetDomain :one
SELECT * FROM domains WHERE name = $1;

-- name: GetDomainByNamespace :one
SELECT * FROM domains
WHERE namespace = $1 AND name = $2;

-- name: GetAllDomainsByNamespace :many
SELECT * FROM domains
WHERE namespace = $1
ORDER BY name;

-- name: GetAllDomains :many
SELECT * FROM domains
ORDER BY name;

-- name: UpdateDomain :one
UPDATE domains
SET status = $1, certificate_expiry = $2, issuer = $3
WHERE name = $4
RETURNING *;

-- name: DeleteDomain :exec
DELETE FROM domains WHERE name = $1;
