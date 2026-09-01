-- name: CreateRecoveryCase :one
INSERT INTO recovery_cases (
    customer_id, subscription_id, payment_id, razorpay_error_code, 
    razorpay_error_desc, amount_at_risk, currency, first_failed_at, recovery_deadline
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW(), NOW() + INTERVAL '14 days'
)
RETURNING id;

-- name: GetRecoveryCaseByID :one
SELECT * FROM recovery_cases WHERE id = $1 LIMIT 1;

-- name: GetActiveCaseBySubscription :one
SELECT id, status FROM recovery_cases
WHERE subscription_id = $1 AND status NOT IN ('RECOVERED', 'UNRECOVERABLE', 'EXPIRED')
LIMIT 1;

-- name: UpdateCaseDiagnosis :exec
UPDATE recovery_cases
SET ai_diagnosis = $2, ai_confidence = $3, decline_category = $4, status = 'DIAGNOSING', updated_at = NOW()
WHERE id = $1;

-- name: UpdateCaseStrategy :exec
UPDATE recovery_cases
SET ai_strategy = $2, status = 'RECOVERING', updated_at = NOW()
WHERE id = $1;

-- name: UpdateCaseStatus :exec
UPDATE recovery_cases
SET status = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateNextRetry :exec
UPDATE recovery_cases
SET next_retry_at = $2, retry_count = retry_count + 1, updated_at = NOW()
WHERE id = $1;

-- name: UpdateCaseRecovered :exec
UPDATE recovery_cases
SET status = 'RECOVERED', amount_recovered = $2, recovered_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: ListRecoveryCases :many
SELECT c.*, 
       a.action_type AS latest_action_type, 
       a.status AS latest_action_status,
       a.channel AS latest_action_channel
FROM recovery_cases c
LEFT JOIN LATERAL (
    SELECT action_type, status, channel
    FROM recovery_actions
    WHERE recovery_case_id = c.id
    ORDER BY created_at DESC
    LIMIT 1
) a ON true
ORDER BY c.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListPendingCases :many
SELECT c.*, 
       a.action_type AS latest_action_type, 
       a.status AS latest_action_status,
       a.channel AS latest_action_channel
FROM recovery_cases c
JOIN LATERAL (
    SELECT action_type, status, channel
    FROM recovery_actions
    WHERE recovery_case_id = c.id
    ORDER BY created_at DESC
    LIMIT 1
) a ON true
WHERE EXISTS (
    SELECT 1 FROM recovery_actions ra 
    WHERE ra.recovery_case_id = c.id AND ra.status = 'PENDING_APPROVAL'
)
ORDER BY c.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetCaseSummary :one
SELECT c.*, 
       a.action_type AS latest_action_type, 
       a.status AS latest_action_status,
       a.channel AS latest_action_channel
FROM recovery_cases c
LEFT JOIN LATERAL (
    SELECT action_type, status, channel
    FROM recovery_actions
    WHERE recovery_case_id = c.id
    ORDER BY created_at DESC
    LIMIT 1
) a ON true
WHERE c.id = $1 LIMIT 1;
