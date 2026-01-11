-- +goose Up
-- +goose StatementBegin
ALTER TABLE email_tokens DROP CONSTRAINT email_tokens_user_id_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE email_tokens ADD CONSTRAINT email_tokens_user_id_key UNIQUE (user_id);
-- +goose StatementEnd
