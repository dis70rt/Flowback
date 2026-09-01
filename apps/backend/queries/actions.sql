-- name: CreateRecoveryAction :one
INSERT INTO recovery_actions (
    recovery_case_id, idempotency_key, action_type, channel, 
    ai_reasoning, discount_percentage, draft_payload, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id;

-- name: GetPendingActions :many
SELECT id, recovery_case_id, action_type, channel, ai_reasoning, discount_percentage, draft_payload
FROM recovery_actions
WHERE status = 'PENDING_APPROVAL'
ORDER BY created_at ASC;

-- name: GetActionsByCase :many
SELECT id, action_type, channel, status, ai_reasoning, draft_payload, discount_percentage, executed_at, created_at
FROM recovery_actions
WHERE recovery_case_id = $1
ORDER BY created_at DESC;

-- name: GetActionByIdempotencyKey :one
SELECT id, status FROM recovery_actions WHERE idempotency_key = $1 LIMIT 1;

-- name: ApproveAction :exec
UPDATE recovery_actions
SET status = 'APPROVED', approved_by_clerk_id = $2, payment_link_id = $3, payment_link_url = $4
WHERE id = $1;

-- name: RejectAction :exec
UPDATE recovery_actions
SET status = 'REJECTED', approved_by_clerk_id = $2
WHERE id = $1;

-- name: MarkActionExecuted :exec
UPDATE recovery_actions
SET status = 'EXECUTED', external_id = $2, executed_at = NOW()
WHERE id = $1;

-- name: UpdateActionAsynqTask :exec
UPDATE recovery_actions
SET asynq_task_id = $2
WHERE id = $1;

-- name: UpdateActionDraft :exec
UPDATE recovery_actions
SET draft_payload = $2, human_edited = TRUE
WHERE id = $1;

-- name: GetActionAndCaseForApproval :one
SELECT a.id as action_id, a.discount_percentage, c.id as case_id, c.amount_at_risk
FROM recovery_actions a
JOIN recovery_cases c ON a.recovery_case_id = c.id
WHERE a.id = $1 LIMIT 1;
