-- +goose Up
ALTER TYPE action_type ADD VALUE 'SEND_CALL';

-- +goose Down
-- Cannot remove a value from ENUM in postgres without recreating it, 
-- but this is a hackathon so we leave it empty or comment it out.
