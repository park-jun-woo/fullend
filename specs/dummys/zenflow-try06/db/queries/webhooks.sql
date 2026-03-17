-- name: WebhookCreate :one
INSERT INTO webhooks (org_id, url, event_type)
VALUES ($1, $2, $3)
RETURNING *;

-- name: WebhookListByOrgID :many
SELECT * FROM webhooks WHERE org_id = $1;

-- name: WebhookListByOrgIDAndEventType :many
SELECT * FROM webhooks WHERE org_id = $1 AND event_type = $2;

-- name: WebhookFindByID :one
SELECT * FROM webhooks WHERE id = $1;

-- name: WebhookFindByIDAndOrgID :one
SELECT * FROM webhooks WHERE id = $1 AND org_id = $2;

-- name: WebhookDelete :exec
DELETE FROM webhooks WHERE id = $1;
