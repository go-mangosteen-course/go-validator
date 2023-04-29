-- name: CreateValidationCode :one
INSERT INTO validation_codes (
  email, code
) VALUES (
  $1, $2
)
RETURNING *;

-- name: CountValidationCodes :one
SELECT count(*) FROM validation_codes WHERE email = $1;
