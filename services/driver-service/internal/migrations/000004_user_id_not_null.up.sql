-- Remove orphaned drivers that have no user_id (cannot be identified or routed)
DELETE FROM drivers WHERE user_id IS NULL;

-- Enforce NOT NULL going forward
ALTER TABLE drivers ALTER COLUMN user_id SET NOT NULL;
