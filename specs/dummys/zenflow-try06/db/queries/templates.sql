-- name: TemplateCreate :one
INSERT INTO templates (source_workflow_id, org_id, title, description, category)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: TemplateFindByID :one
SELECT * FROM templates WHERE id = $1;

-- name: TemplateFindBySourceWorkflowID :one
SELECT * FROM templates WHERE source_workflow_id = $1;

-- name: TemplateList :many
SELECT * FROM templates;

-- name: TemplateIncrementCloneCount :exec
UPDATE templates SET clone_count = clone_count + 1 WHERE id = $1;
