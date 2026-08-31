-- +goose Up

-- 1. Add Location for the Web Search Tool
ALTER TABLE customers 
ADD COLUMN city VARCHAR(100),
ADD COLUMN state VARCHAR(100);

-- 2. Add HITL Statuses
ALTER TYPE action_status ADD VALUE 'PENDING_APPROVAL';
ALTER TYPE action_status ADD VALUE 'APPROVED';
ALTER TYPE action_status ADD VALUE 'REJECTED';

-- 3. Add UI Drafts, Discounts, Queue Tracking, and Clerk Auth
ALTER TABLE recovery_actions
ADD COLUMN discount_percentage INT DEFAULT 0,
ADD COLUMN draft_subject TEXT,
ADD COLUMN draft_body TEXT,
ADD COLUMN asynq_task_id VARCHAR(255),
ADD COLUMN approved_by_clerk_id VARCHAR(255);

-- +goose Down
ALTER TABLE customers 
DROP COLUMN city,
DROP COLUMN state;

ALTER TABLE recovery_actions
DROP COLUMN discount_percentage,
DROP COLUMN draft_subject,
DROP COLUMN draft_body,
DROP COLUMN asynq_task_id,
DROP COLUMN approved_by_clerk_id;
