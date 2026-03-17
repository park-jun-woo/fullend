-- name: ProposalCreate :one
INSERT INTO proposals (gig_id, freelancer_id, bid_amount, status)
VALUES ($1, $2, $3, $4)
RETURNING id, gig_id, freelancer_id, bid_amount, status;

-- name: ProposalFindByID :one
SELECT id, gig_id, freelancer_id, bid_amount, status FROM proposals WHERE id = $1;

-- name: ProposalUpdateStatus :exec
UPDATE proposals SET status = $2 WHERE id = $1;
