-- name: WorkflowCreate :one
INSERT INTO workflows (org_id, title, trigger_event, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: WorkflowFindByID :one
SELECT * FROM workflows WHERE id = $1;

-- name: WorkflowFindByIDAndOrgID :one
SELECT * FROM workflows WHERE id = $1 AND org_id = $2;

-- name: WorkflowListByOrgID :many
SELECT * FROM workflows WHERE org_id = $1;

-- name: WorkflowUpdateStatus :exec
UPDATE workflows SET status = $2 WHERE id = $1;

-- name: WorkflowCreateVersion :one
INSERT INTO workflows (org_id, title, trigger_event, status, version, root_workflow_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: WorkflowListVersions :many
SELECT * FROM workflows WHERE (root_workflow_id = $1 OR id = $1) AND org_id = $2 ORDER BY version DESC;
