-- +goose Up
-- 1. Action Metadata
ALTER TABLE recovery_actions 
ADD COLUMN human_edited BOOLEAN NOT NULL DEFAULT FALSE;

-- 2. Smart Trigger: Auto-close cases on payment success
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION process_payment_webhook()
RETURNS TRIGGER AS $$
DECLARE
    case_id_val UUID;
    cust_id_val UUID;
    amount_val  BIGINT;
BEGIN
    IF NEW.event_type IN ('payment_link.paid', 'payment.captured') THEN
        -- payment_link.paid path: payload.payment_link.entity.notes.recovery_case_id
        -- payment.captured path:  payload.payment.entity.notes.recovery_case_id
        -- Try both paths; whichever is non-null wins.
        case_id_val := COALESCE(
            (NEW.payload #>> '{payload,payment_link,entity,notes,recovery_case_id}')::UUID,
            (NEW.payload #>> '{payload,payment,entity,notes,recovery_case_id}')::UUID
        );

        IF case_id_val IS NOT NULL THEN
            UPDATE recovery_cases
            SET
                status           = 'RECOVERED',
                recovered_at     = NOW(),
                -- Set amount_recovered so metrics queries return real numbers
                amount_recovered = amount_at_risk
            WHERE id = case_id_val
            RETURNING customer_id, amount_at_risk INTO cust_id_val, amount_val;

            IF cust_id_val IS NOT NULL THEN
                UPDATE customers
                SET successful_payments = successful_payments + 1,
                    total_payments      = total_payments + 1
                WHERE id = cust_id_val;
            END IF;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_process_webhook
AFTER INSERT ON webhook_events
FOR EACH ROW EXECUTE FUNCTION process_payment_webhook();

-- +goose Down
DROP TRIGGER IF EXISTS trg_process_webhook ON webhook_events;
DROP FUNCTION IF EXISTS process_payment_webhook();

ALTER TABLE recovery_actions 
DROP COLUMN IF EXISTS human_edited;
