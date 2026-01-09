-- +goose Up
-- +goose StatementBegin
CREATE TABLE email_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    token_hash text NOT NULL UNIQUE,
    type TEXT NOT NULL,
    expires_at timestamp NOT NULL,
    used_at timestamp,
    created_at timestamp NOT NULL DEFAULT now()
);

CREATE INDEX idx_email_token_user_type ON email_tokens(user_id, type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_email_token_user_type;
DROP TABLE IF EXISTS email_tokens;
-- +goose StatementEnd
