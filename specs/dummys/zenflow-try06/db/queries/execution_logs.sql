-- name: ExecutionLogCreate :one
INSERT INTO execution_logs (workflow_id, org_id, status, credits_spent)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ExecutionLogCreateWithReport :one
INSERT INTO execution_logs (workflow_id, org_id, status, credits_spent, report_key)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ExecutionLogFindByIDAndOrgID :one
SELECT * FROM execution_logs WHERE id = $1 AND org_id = $2;

-- name: ExecutionLogListByWorkflowID :many
SELECT * FROM execution_logs WHERE workflow_id = $1 ORDER BY executed_at DESC;
