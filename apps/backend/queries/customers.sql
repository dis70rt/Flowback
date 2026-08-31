-- name: GetCustomerProfile :one
SELECT 
    id, razorpay_customer_id, name, value_tier, tenure, 
    preferred_channel, city, state, failed_payments, reliability_score
FROM customers
WHERE razorpay_customer_id = $1 LIMIT 1;

-- name: IncrementFailedPayment :exec
UPDATE customers
SET failed_payments = failed_payments + 1, updated_at = NOW()
WHERE id = $1;

-- name: IncrementSuccessfulPayment :exec
UPDATE customers
SET successful_payments = successful_payments + 1, total_payments = total_payments + 1, updated_at = NOW()
WHERE id = $1;

-- name: InsertCustomer :one
INSERT INTO customers (
    razorpay_customer_id, email, phone, name, value_tier, tenure, preferred_channel, city, state
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
ON CONFLICT (razorpay_customer_id) DO UPDATE 
SET email = EXCLUDED.email, phone = EXCLUDED.phone, name = EXCLUDED.name, updated_at = NOW()
RETURNING id;

-- name: GetCustomerByEmailOrPhone :one
SELECT id
FROM customers
WHERE (email = $1 AND email != '') OR (phone = $2 AND phone != '')
LIMIT 1;
