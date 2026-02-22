ALTER TABLE drivers ADD COLUMN IF NOT EXISTS car_plate VARCHAR(50);
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS package_slug VARCHAR(50);

UPDATE drivers d
SET car_plate = c.car_plate,
    package_slug = c.package_slug
FROM cars c
WHERE d.current_car_id = c.id;

ALTER TABLE drivers DROP COLUMN IF EXISTS current_car_id;
DROP TABLE IF EXISTS cars;
CREATE INDEX IF NOT EXISTS idx_drivers_package_slug ON drivers (package_slug);
