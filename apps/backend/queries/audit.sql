-- name: InsertAuditLog :exec
INSERT INTO audit_logs (source_binary, event_type, actor, trace_id, entity_id, payload)
VALUES ($1, $2, $3, $4, $5, $6);
