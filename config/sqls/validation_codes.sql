-- name: CreateValidationCode :one
INSERT INTO validation_codes (
  email, code
) VALUES (
  $1, $2
)
RETURNING *;

-- name: CountValidationCodes :one
SELECT count(*) FROM validation_codes WHERE email = $1;

-- name: FindValidationCode :one
SELECT * FROM validation_codes
WHERE
  email = $1
  AND
  code = $2
  AND
  used_at = null
ORDER BY created_at desc
LIMIT 1;
