-- name: InsertWebhookEvent :one
INSERT INTO webhook_events (razorpay_event_id, event_type, payload, signature)
VALUES ($1, $2, $3, $4)
ON CONFLICT (razorpay_event_id) DO NOTHING
RETURNING id;

-- name: MarkWebhookProcessed :exec
UPDATE webhook_events
SET processed = TRUE, processed_at = NOW()
WHERE id = $1;
