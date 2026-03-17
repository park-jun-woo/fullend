-- name: GigCreate :one
INSERT INTO gigs (client_id, title, description, budget, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, client_id, title, description, budget, status, freelancer_id, created_at;

-- name: GigFindByID :one
SELECT id, client_id, title, description, budget, status, freelancer_id, created_at FROM gigs WHERE id = $1;

-- name: GigList :many
SELECT id, client_id, title, description, budget, status, freelancer_id, created_at FROM gigs;

-- name: GigUpdateStatus :exec
UPDATE gigs SET status = $2 WHERE id = $1;

-- name: GigUpdateFreelancer :exec
UPDATE gigs SET freelancer_id = $2 WHERE id = $1;
