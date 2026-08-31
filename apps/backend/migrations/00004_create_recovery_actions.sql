-- +goose Up
CREATE TYPE action_type AS ENUM (
    'SILENT_RETRY', 'SEND_PAYMENT_LINK', 'SEND_EMAIL', 'SEND_SMS',
    'SEND_WHATSAPP', 'ESCALATE_TO_HUMAN', 'MARK_UNRECOVERABLE'
);

CREATE TYPE action_status AS ENUM (
    'PENDING', 'EXECUTED', 'FAILED', 'BLOCKED_BY_POLICY'
);

CREATE TABLE recovery_actions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recovery_case_id    UUID NOT NULL REFERENCES recovery_cases(id),
    idempotency_key     VARCHAR(255) UNIQUE NOT NULL,
    action_type         action_type NOT NULL,
    channel             VARCHAR(20),
    ai_reasoning        TEXT,
    ai_confidence       FLOAT,
    policy_approved     BOOLEAN NOT NULL DEFAULT FALSE,
    policy_block_reason TEXT,
    status              action_status NOT NULL DEFAULT 'PENDING',
    external_id         VARCHAR(255),
    error_message       TEXT,
    payment_link_id     VARCHAR(255),
    payment_link_url    TEXT,
    payment_link_expires_at TIMESTAMPTZ,
    executed_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_recovery_actions_case_id ON recovery_actions(recovery_case_id);
CREATE INDEX idx_recovery_actions_idempotency ON recovery_actions(idempotency_key);

-- +goose Down
DROP TABLE IF EXISTS recovery_actions;
DROP TYPE IF EXISTS action_status;
DROP TYPE IF EXISTS action_type;
