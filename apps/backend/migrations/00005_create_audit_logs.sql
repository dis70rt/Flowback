-- +goose Up
CREATE TABLE audit_logs (
    id              BIGSERIAL PRIMARY KEY,
    source_binary   VARCHAR(50) NOT NULL,
    event_type      VARCHAR(100) NOT NULL,
    actor           VARCHAR(100),
    trace_id        UUID,
    entity_id       VARCHAR(255),
    payload         JSONB NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_event_type ON audit_logs(event_type);
CREATE INDEX idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX idx_audit_logs_payload_gin ON audit_logs USING GIN(payload);

-- +goose Down
DROP TABLE IF EXISTS audit_logs;
