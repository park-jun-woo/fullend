-- name: OrganizationFindByID :one
SELECT * FROM organizations WHERE id = $1;

-- name: OrganizationDeductCredit :exec
UPDATE organizations SET credits_balance = credits_balance - 1 WHERE id = $1;
