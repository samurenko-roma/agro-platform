CREATE TABLE weather_locations
(
    id                UUID         PRIMARY KEY,
    farm_id           UUID         NOT NULL,
    production_unit_id UUID        NULL,

    name              VARCHAR(255) NOT NULL,
    latitude          DECIMAL(10, 6) NOT NULL,
    longitude         DECIMAL(10, 6) NOT NULL,
    timezone          VARCHAR(100) NOT NULL DEFAULT 'UTC',
    is_default        BOOLEAN      NOT NULL DEFAULT false,

    created_at        TIMESTAMPTZ  NOT NULL,
    updated_at        TIMESTAMPTZ  NOT NULL,
    archived_at       TIMESTAMPTZ  NULL
);

CREATE INDEX idx_weather_locations_farm_id ON weather_locations (farm_id);

-- Только одна дефолтная локация на ферму
CREATE UNIQUE INDEX idx_weather_locations_default_per_farm
    ON weather_locations (farm_id)
    WHERE is_default = true AND archived_at IS NULL;
