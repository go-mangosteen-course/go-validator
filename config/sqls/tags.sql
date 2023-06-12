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

-- name: UpdateTag :one
UPDATE tags
SET
  user_id = @user_id,
  name = CASE WHEN @name::varchar = '' THEN name ELSE @name END,
  sign = CASE WHEN @sign::varchar = '' THEN sign ELSE @sign END,
  kind = CASE WHEN @kind::varchar = '' THEN kind ELSE @kind END
WHERE id = @id
RETURNING *;

-- name: DeleteTag :exec
UPDATE tags
SET deleted_at = now()
WHERE id = @id;

-- name: FindTag :one
SELECT * FROM tags
WHERE id = @id AND deleted_at IS NULL;
-- name: ListTags :many
SELECT * FROM tags
WHERE
kind = @kind AND user_id = @user_id AND deleted_at IS NULL
ORDER BY created_at DESC
OFFSET $1
LIMIT $2;

