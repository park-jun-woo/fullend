-- name: UserCreate :one
INSERT INTO users (email, password_hash, role, name)
VALUES ($1, $2, $3, $4)
RETURNING id, email, password_hash, role, name;

-- name: UserFindByID :one
SELECT id, email, password_hash, role, name FROM users WHERE id = $1;

-- name: UserFindByEmail :one
SELECT id, email, password_hash, role, name FROM users WHERE email = $1;
