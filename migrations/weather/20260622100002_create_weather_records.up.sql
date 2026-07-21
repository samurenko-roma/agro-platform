CREATE TABLE weather_records
(
    id          UUID         PRIMARY KEY,
    location_id UUID         NOT NULL REFERENCES weather_locations (id) ON DELETE CASCADE,
    farm_id     UUID         NOT NULL,

    kind        VARCHAR(20)  NOT NULL, -- CURRENT | FORECAST | HISTORICAL
    source      VARCHAR(50)  NOT NULL, -- OPEN_METEO | SENSOR | CUSTOM

    timestamp   TIMESTAMPTZ  NOT NULL, -- время наблюдения / прогноза
    forecast_at TIMESTAMPTZ  NULL,     -- когда сделан прогноз (только для FORECAST)

    data        JSONB        NOT NULL DEFAULT '{}',

    created_at  TIMESTAMPTZ  NOT NULL
);

CREATE INDEX idx_weather_records_location_kind     ON weather_records (location_id, kind);
CREATE INDEX idx_weather_records_location_timestamp ON weather_records (location_id, timestamp DESC);
CREATE INDEX idx_weather_records_farm_kind         ON weather_records (farm_id, kind);

-- Не храним дубли: один источник — одна запись на момент времени
CREATE UNIQUE INDEX idx_weather_records_unique_point
    ON weather_records (location_id, kind, source, timestamp);
