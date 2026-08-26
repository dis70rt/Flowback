-- +goose Up
CREATE TABLE webhook_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    razorpay_event_id VARCHAR(255) UNIQUE NOT NULL,
    event_type      VARCHAR(100) NOT NULL,
    payload         JSONB NOT NULL,
    signature       VARCHAR(255) NOT NULL,
    processed       BOOLEAN NOT NULL DEFAULT FALSE,
    processed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_webhook_events_razorpay_id ON webhook_events(razorpay_event_id);
CREATE INDEX idx_webhook_events_processed ON webhook_events(processed) WHERE processed = FALSE;

-- +goose Down
DROP TABLE IF EXISTS webhook_events;
