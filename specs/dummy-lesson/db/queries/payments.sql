-- name: FindByID :one
SELECT * FROM payments WHERE id = $1;

-- name: ListByUser :many
SELECT * FROM payments WHERE user_id = $1 ORDER BY created_at DESC;

-- name: Create :one
INSERT INTO payments (user_id, enrollment_id, amount, method, status)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateStatus :exec
UPDATE payments SET status = $2 WHERE id = $1;
