CREATE TABLE
   IF NOT EXISTS drivers (
      id UUID PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      profile_picture TEXT,
      car_plate VARCHAR(50),
      geohash VARCHAR(20),
      package_slug VARCHAR(50) NOT NULL,
      latitude DECIMAL(10, 8),
      longitude DECIMAL(11, 8),
      is_available BOOLEAN DEFAULT true,
      created_at TIMESTAMP
      WITH
         TIME ZONE DEFAULT NOW (),
         updated_at TIMESTAMP
      WITH
         TIME ZONE DEFAULT NOW ()
   );

CREATE INDEX idx_drivers_package_slug ON drivers (package_slug);

CREATE INDEX idx_drivers_is_available ON drivers (is_available);

CREATE INDEX idx_drivers_geohash ON drivers (geohash);