-- name: GigCreate :one
INSERT INTO gigs (client_id, title, description, budget, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GigFindByID :one
SELECT * FROM gigs WHERE id = $1;

-- name: GigList :many
SELECT * FROM gigs;

-- name: GigUpdateStatus :exec
UPDATE gigs SET status = $2 WHERE id = $1;

-- name: GigUpdateStatusAndFreelancer :exec
UPDATE gigs SET status = $2, freelancer_id = $3 WHERE id = $1;
