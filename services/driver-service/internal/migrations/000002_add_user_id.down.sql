ALTER TABLE drivers DROP COLUMN IF EXISTS user_id;
DROP INDEX IF EXISTS idx_drivers_user_id;
