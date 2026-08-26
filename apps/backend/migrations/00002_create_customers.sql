-- +goose Up
CREATE TABLE customers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    razorpay_customer_id VARCHAR(255) UNIQUE,
    email           VARCHAR(255),
    phone           VARCHAR(20),
    name            VARCHAR(255),
    value_tier      VARCHAR(20) DEFAULT 'MEDIUM' CHECK (value_tier IN ('HIGH', 'MEDIUM', 'LOW')),
    tenure          VARCHAR(20) DEFAULT 'NEW' CHECK (tenure IN ('NEW', 'ESTABLISHED', 'LOYAL')),
    preferred_channel VARCHAR(20) DEFAULT 'EMAIL' CHECK (preferred_channel IN ('EMAIL', 'SMS', 'WHATSAPP')),
    total_payments       INT NOT NULL DEFAULT 0,
    successful_payments  INT NOT NULL DEFAULT 0,
    failed_payments      INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_customers_razorpay_id ON customers(razorpay_customer_id);
CREATE INDEX idx_customers_email ON customers(email);

-- +goose Down
DROP TABLE IF EXISTS customers;
