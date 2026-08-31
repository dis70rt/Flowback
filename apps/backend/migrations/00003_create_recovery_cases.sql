-- +goose Up
CREATE TYPE recovery_status AS ENUM (
    'DETECTED', 'DIAGNOSING', 'RECOVERING', 'RECOVERED',
    'ESCALATED', 'UNRECOVERABLE', 'EXPIRED'
);

CREATE TYPE decline_type AS ENUM (
    'SOFT', 'HARD', 'AMBIGUOUS'
);

CREATE TABLE recovery_cases (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id         UUID REFERENCES customers(id),
    subscription_id     VARCHAR(255) NOT NULL,
    payment_id          VARCHAR(255),
    razorpay_error_code VARCHAR(255),
    razorpay_error_desc TEXT,
    decline_category    decline_type,
    ai_diagnosis        TEXT,
    ai_confidence       FLOAT,
    status              recovery_status NOT NULL DEFAULT 'DETECTED',
    retry_count         INT NOT NULL DEFAULT 0,
    max_retries         INT NOT NULL DEFAULT 6,
    contact_count       INT NOT NULL DEFAULT 0,
    max_contacts        INT NOT NULL DEFAULT 4,
    amount_at_risk      BIGINT NOT NULL,
    currency            VARCHAR(10) NOT NULL DEFAULT 'INR',
    amount_recovered    BIGINT,
    first_failed_at     TIMESTAMPTZ NOT NULL,
    recovery_deadline   TIMESTAMPTZ NOT NULL,
    recovered_at        TIMESTAMPTZ,
    next_retry_at       TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_recovery_cases_subscription_id ON recovery_cases(subscription_id);
CREATE INDEX idx_recovery_cases_status ON recovery_cases(status);
CREATE INDEX idx_recovery_cases_next_retry ON recovery_cases(next_retry_at) WHERE status = 'RECOVERING';

-- +goose Down
DROP TABLE IF EXISTS recovery_cases;
DROP TYPE IF EXISTS decline_type;
DROP TYPE IF EXISTS recovery_status;
