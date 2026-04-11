-- name: ExecutionLogCreate :one
INSERT INTO execution_logs (workflow_id, org_id, status, credits_spent) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: ExecutionLogListByWorkflow :many
SELECT * FROM execution_logs WHERE workflow_id = $1 ORDER BY executed_at DESC;
