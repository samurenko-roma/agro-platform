CREATE TABLE spatial_layout_snapshots
(
    id          UUID        PRIMARY KEY,
    farm_id     UUID        NOT NULL,
    version     INTEGER     NOT NULL,
    description TEXT        NULL,
    created_by  UUID        NOT NULL,
    units       JSONB       NOT NULL DEFAULT '[]',
    created_at  TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_spatial_layout_snapshots_farm_id ON spatial_layout_snapshots (farm_id);

-- Версии снапшотов одной фермы не должны повторяться — GetLatest полагается
-- на монотонный ORDER BY version DESC.
CREATE UNIQUE INDEX ux_spatial_layout_snapshots_farm_version
    ON spatial_layout_snapshots (farm_id, version);