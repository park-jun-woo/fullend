-- name: FindByID :one
SELECT * FROM lessons WHERE id = $1;

-- name: ListByCourse :many
SELECT * FROM lessons WHERE course_id = $1 ORDER BY sort_order ASC;

-- name: Create :one
INSERT INTO lessons (course_id, title, video_url, sort_order)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: Update :exec
UPDATE lessons SET title = $2, video_url = $3, sort_order = $4
WHERE id = $1;

-- name: Delete :exec
DELETE FROM lessons WHERE id = $1;
