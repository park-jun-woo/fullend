-- name: ProposalFindByID :one
SELECT * FROM proposals WHERE id = $1;

-- name: ProposalCreate :one
INSERT INTO proposals (gig_id, freelancer_id, bid_amount, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ProposalUpdateStatus :exec
UPDATE proposals SET status = $1 WHERE id = $2;
