-- name: FindByID :one
SELECT * FROM reviews WHERE id = $1;

-- name: FindByCourseAndUser :one
SELECT * FROM reviews WHERE course_id = $1 AND user_id = $2;

-- name: ListByCourse :many
SELECT * FROM reviews WHERE course_id = $1 ORDER BY created_at DESC;

-- name: Create :one
INSERT INTO reviews (user_id, course_id, rating, comment)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: Delete :exec
DELETE FROM reviews WHERE id = $1;
