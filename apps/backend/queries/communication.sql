-- name: GetCommunicationHistory :many
SELECT *
FROM communication_history
WHERE customer_id = (SELECT id FROM customers WHERE razorpay_customer_id = $1 LIMIT 1)
ORDER BY sent_at DESC LIMIT 5;

-- name: InsertCommunicationHistory :one
INSERT INTO communication_history (recovery_case_id, customer_id, channel, message_sid)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: UpdateCommunicationStatus :exec
UPDATE communication_history
SET status = $2, delivered_at = CASE WHEN $2 = 'DELIVERED' THEN NOW() ELSE delivered_at END,
opened_at = CASE WHEN $2 = 'OPENED' THEN NOW() ELSE opened_at END,
clicked_at = CASE WHEN $2 = 'CLICKED' THEN NOW() ELSE clicked_at END
WHERE message_sid = $1;
