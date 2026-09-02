-- name: GetDashboardOverview :one
SELECT 
    COALESCE(SUM(amount_at_risk), 0)::BIGINT as total_amount_at_risk,
    COALESCE(SUM(amount_recovered), 0)::BIGINT as total_amount_recovered,
    COUNT(CASE WHEN status = 'RECOVERED' THEN 1 END)::INT as recovered_cases_count,
    COUNT(CASE WHEN status IN ('DETECTED', 'DIAGNOSING', 'RECOVERING') THEN 1 END)::INT as active_cases_count,
    COUNT(CASE WHEN status IN ('UNRECOVERABLE', 'EXPIRED', 'ESCALATED') THEN 1 END)::INT as failed_cases_count
FROM recovery_cases;

-- name: GetRecoveryTrends :many
SELECT 
    DATE_TRUNC('day', created_at)::DATE as date,
    COALESCE(SUM(amount_at_risk), 0)::BIGINT as daily_failed,
    COALESCE(SUM(amount_recovered), 0)::BIGINT as daily_recovered
FROM recovery_cases
WHERE created_at >= NOW() - INTERVAL '30 days'
GROUP BY DATE_TRUNC('day', created_at)
ORDER BY date ASC;

-- name: GetChannelDistribution :many
SELECT 
    channel,
    COUNT(*)::INT as count
FROM recovery_actions
WHERE channel IS NOT NULL
GROUP BY channel
ORDER BY count DESC;

-- name: GetPipelineStatus :many
SELECT 
    status,
    COUNT(*)::INT as count
FROM recovery_actions
GROUP BY status
ORDER BY count DESC;
