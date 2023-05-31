-- name: CreateTag :one
INSERT INTO tags (
  user_id,
  name,
  sign,
  kind
) VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING *;
