-- +goose Up

-- 1. Communication History (for Dynamic Channel Escalation tool)
--    Tracks every outbound message and its engagement.
--    opened_at/clicked_at only populated when the channel supports it:
--    WhatsApp -> all three (blue ticks), Email -> clicked_at only (pixel unreliable), SMS -> NULL always
CREATE TABLE communication_history (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recovery_case_id    UUID NOT NULL REFERENCES recovery_cases(id),
    customer_id         UUID NOT NULL REFERENCES customers(id),
    channel             VARCHAR(20) NOT NULL CHECK (channel IN ('EMAIL', 'SMS', 'WHATSAPP')),
    status              VARCHAR(30) NOT NULL DEFAULT 'SENT'
                            CHECK (status IN ('SENT', 'DELIVERED', 'OPENED', 'CLICKED', 'FAILED')),
    message_sid         VARCHAR(255),
    sent_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    delivered_at        TIMESTAMPTZ,
    opened_at           TIMESTAMPTZ,
    clicked_at          TIMESTAMPTZ
);
CREATE INDEX idx_comm_history_customer    ON communication_history(customer_id);
CREATE INDEX idx_comm_history_case        ON communication_history(recovery_case_id);
CREATE INDEX idx_comm_history_message_sid ON communication_history(message_sid) WHERE message_sid IS NOT NULL;

-- 2. Store the full AI Strategy JSON on the recovery_case
--    Lets the UI re-render the "AI Reasoning" panel on page refresh
ALTER TABLE recovery_cases
ADD COLUMN ai_strategy JSONB;

-- 3. Fast partial index for the Supervisor UI main query
CREATE INDEX idx_recovery_actions_pending
    ON recovery_actions(created_at DESC)
    WHERE status = 'PENDING_APPROVAL';

-- 4. Customer reliability score (used by StrategyAgent tool)
ALTER TABLE customers
ADD COLUMN reliability_score FLOAT NOT NULL DEFAULT 1.0;

-- +goose Down
DROP TABLE IF EXISTS communication_history;
DROP INDEX IF EXISTS idx_recovery_actions_pending;
ALTER TABLE recovery_cases DROP COLUMN IF EXISTS ai_strategy;
ALTER TABLE customers DROP COLUMN IF EXISTS reliability_score;
