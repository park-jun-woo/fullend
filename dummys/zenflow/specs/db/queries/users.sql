-- name: UserCreate :one
INSERT INTO users (org_id, email, password_hash, role, name) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UserFindByID :one
SELECT * FROM users WHERE id = $1;

-- name: UserFindByEmail :one
SELECT * FROM users WHERE email = $1;
