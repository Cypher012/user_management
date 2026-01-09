-- name: CreateEmailToken :exec
INSERT INTO email_tokens (
    user_id,
    token_hash,
    type,
    expires_at
) VALUES ($1, $2, $3, $4);

-- name: GetValidEmailToken :one
SELECT *
FROM email_tokens
WHERE token_hash = $1
  AND type = $2
  AND used_at IS NULL
  AND expires_at > now();

-- name: MarkEmailTokenUsed :exec
UPDATE email_tokens
SET used_at = now()
WHERE id = $1;
