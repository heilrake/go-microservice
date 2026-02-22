ALTER TABLE drivers ADD COLUMN IF NOT EXISTS user_id UUID;
CREATE UNIQUE INDEX IF NOT EXISTS idx_drivers_user_id ON drivers (user_id);
