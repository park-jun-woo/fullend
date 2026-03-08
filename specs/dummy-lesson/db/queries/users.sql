-- name: FindByID :one
SELECT * FROM users WHERE id = $1;

-- name: FindByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: Create :one
INSERT INTO users (email, password_hash, name, role)
VALUES ($1, $2, $3, $4) RETURNING *;
