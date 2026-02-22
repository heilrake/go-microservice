-- Create cars table
CREATE TABLE IF NOT EXISTS cars (
    id UUID PRIMARY KEY,
    driver_id UUID NOT NULL REFERENCES drivers(id) ON DELETE CASCADE,
    car_plate VARCHAR(50) NOT NULL,
    package_slug VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_cars_driver_id ON cars (driver_id);
CREATE INDEX idx_cars_package_slug ON cars (package_slug);

-- Add current_car_id to drivers (which car they're using when available)
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS current_car_id UUID REFERENCES cars(id);
CREATE INDEX idx_drivers_current_car ON drivers (current_car_id);

-- Migrate existing drivers: create car from car_plate/package_slug, link it
INSERT INTO cars (id, driver_id, car_plate, package_slug)
SELECT gen_random_uuid(), id, COALESCE(car_plate, 'MIGRATED'), package_slug
FROM drivers
WHERE package_slug IS NOT NULL AND package_slug != '';

UPDATE drivers d
SET current_car_id = (SELECT id FROM cars c WHERE c.driver_id = d.id LIMIT 1)
WHERE EXISTS (SELECT 1 FROM cars c WHERE c.driver_id = d.id);

-- Drop old columns from drivers
ALTER TABLE drivers DROP COLUMN IF EXISTS car_plate;
ALTER TABLE drivers DROP COLUMN IF EXISTS package_slug;

-- Drop old index
DROP INDEX IF EXISTS idx_drivers_package_slug;
