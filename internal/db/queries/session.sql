-- name: GetSessionById :one
SELECT *
FROM sessions
WHERE id = $1
    AND revoked_at IS NULL
    AND expires_at > now();

-- name: CreateSession :one
INSERT INTO sessions
(user_id, refresh_token_hash, device_name, device_info, ip_address, expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListActiveUserSessions :many
SELECT *
FROM sessions
WHERE user_id = $1
  AND revoked_at IS NULL
  AND expires_at > now()
ORDER BY last_used DESC;

-- name: RevokeSession :exec
UPDATE sessions
SET revoked_at = now()
WHERE id = $1
  AND revoked_at IS NULL;

-- name: GetSessionByRefreshToken :one
SELECT *
FROM sessions
WHERE refresh_token_hash = $1
  AND revoked_at IS NULL
  AND expires_at > now();

-- name: UpdateSessionLastUsed :exec
UPDATE sessions
SET last_used = now()
WHERE id = $1
  AND revoked_at IS NULL
  AND expires_at > now();

-- name: RotateSessionToken :exec
UPDATE sessions
SET refresh_token_hash = $2,
    last_used = now()
WHERE id = $1
  AND revoked_at IS NULL
  AND expires_at > now();

-- name: RevokeAllUserSessions :exec
UPDATE sessions
SET revoked_at = now()
WHERE user_id = $1
  AND revoked_at IS NULL;

-- name: RevokeAllSessions :exec
UPDATE sessions
SET revoked_at = now()
WHERE revoked_at IS NULL;
