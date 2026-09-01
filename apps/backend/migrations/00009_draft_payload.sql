-- +goose Up
ALTER TABLE recovery_actions
DROP COLUMN IF EXISTS draft_subject,
DROP COLUMN IF EXISTS draft_body,
ADD COLUMN draft_payload JSONB;

-- +goose Down
ALTER TABLE recovery_actions
DROP COLUMN IF EXISTS draft_payload,
ADD COLUMN draft_subject TEXT,
ADD COLUMN draft_body TEXT;
