-- name: LogWebhookEvent :exec
INSERT INTO webhook_events (razorpay_event_id, event_type, payload, signature)
VALUES ($1, $2, $3, $4);
