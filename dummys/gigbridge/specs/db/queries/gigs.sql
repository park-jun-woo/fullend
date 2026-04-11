-- name: GigCreate :one
INSERT INTO gigs (client_id, title, description, budget) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GigFindByID :one
SELECT * FROM gigs WHERE id = $1;

-- name: GigList :many
SELECT * FROM gigs;

-- name: GigUpdateStatus :exec
UPDATE gigs SET status = $2 WHERE id = $1;

-- name: GigAssignFreelancer :exec
UPDATE gigs SET freelancer_id = $2, status = 'in_progress' WHERE id = $1;
