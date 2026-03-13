-- name: OrganizationCreate :one
INSERT INTO organizations (name, plan_type, credits_balance)
VALUES ($1, $2, $3)
RETURNING *;

-- name: OrganizationFindByID :one
SELECT * FROM organizations WHERE id = $1;

-- name: OrganizationUpdateCredits :exec
UPDATE organizations SET credits_balance = credits_balance - $2 WHERE id = $1;
