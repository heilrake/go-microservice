-- Fix: is_available should default to false (driver starts as offline)
-- Also reset any drivers with is_available=true but no current_car_id (zombie state from GORM default bug)
ALTER TABLE drivers ALTER COLUMN is_available SET DEFAULT false;

UPDATE drivers SET is_available = false WHERE current_car_id IS NULL AND is_available = true;
