-- name: ActionCreate :one
INSERT INTO actions (workflow_id, action_type, payload_template, sequence_order)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ActionListByWorkflowID :many
SELECT * FROM actions WHERE workflow_id = $1 ORDER BY sequence_order ASC;

-- name: ActionCopyToWorkflow :exec
INSERT INTO actions (workflow_id, action_type, payload_template, sequence_order)
SELECT $1, action_type, payload_template, sequence_order FROM actions AS src WHERE src.workflow_id = $2;
