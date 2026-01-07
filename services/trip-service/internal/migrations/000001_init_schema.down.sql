-- Drop trips first (depends on ride_fares)
DROP INDEX IF EXISTS idx_trips_status;

DROP INDEX IF EXISTS idx_trips_user_id;

DROP TABLE IF EXISTS trips;

-- Then drop ride_fares
DROP INDEX IF EXISTS idx_ride_fares_user_id;

DROP TABLE IF EXISTS ride_fares;