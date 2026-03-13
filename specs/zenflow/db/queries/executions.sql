-- name: ExecutionCreate :one
INSERT INTO executions (workflow_id, org_id, status, credits_spent)
VALUES ($1, $2, $3, $4)
RETURNING *;
