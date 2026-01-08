-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING id, email, is_verified, is_active, created_at;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: SetUserEmailVerified :exec
UPDATE users
SET is_verified = true,
    updated_at = now()
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2,
    updated_at = now()
WHERE id = $1;

-- name: EmailExists :one
SELECT EXISTS (
  SELECT 1 FROM users WHERE email = $1
);
