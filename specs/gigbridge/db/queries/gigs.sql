-- name: GigFindByID :one
SELECT * FROM gigs WHERE id = $1;

-- name: GigList :many
SELECT * FROM gigs;

-- name: GigCreate :one
INSERT INTO gigs (client_id, title, description, budget, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GigUpdateStatus :exec
UPDATE gigs SET status = $1 WHERE id = $2;

-- name: GigAssignFreelancer :exec
UPDATE gigs SET freelancer_id = $1 WHERE id = $2;
