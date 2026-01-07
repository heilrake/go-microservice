-- Create ride_fares table first (trips depends on it)
CREATE TABLE IF NOT EXISTS ride_fares (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    package_slug VARCHAR(50) NOT NULL,
    total_price_in_cents DECIMAL(12, 2) NOT NULL,
    route_data JSONB,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_ride_fares_user_id ON ride_fares (user_id);

-- Create trips table with foreign key to ride_fares
CREATE TABLE IF NOT EXISTS trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    ride_fare_id UUID REFERENCES ride_fares(id) ON DELETE SET NULL,
    driver_id VARCHAR(255),
    driver_name VARCHAR(255),
    driver_car_plate VARCHAR(50),
    driver_profile_picture TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_trips_user_id ON trips (user_id);
CREATE INDEX idx_trips_status ON trips (status);

