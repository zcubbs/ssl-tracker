-- name: InsertDomain :one
INSERT INTO domains
    (name,status,owner)
VALUES ($1,$2,$3)
RETURNING *;

-- name: GetDomain :one
SELECT * FROM domains WHERE name = $1;

-- name: GetDomainByOwner :one
SELECT * FROM domains
WHERE owner = $1 AND name = $2;

-- name: GetAllDomainsByOwner :many
SELECT * FROM domains
WHERE owner = $1
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
