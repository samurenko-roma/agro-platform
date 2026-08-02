ALTER TABLE production_allocations   ADD COLUMN farm_id UUID;
ALTER TABLE production_plantings     ADD COLUMN farm_id UUID;
ALTER TABLE production_harvest_batch ADD COLUMN farm_id UUID;

-- Backfill по уже существующим данным через production_growing_cycles,
-- у которого farm_id уже есть с самого начала.
UPDATE production_allocations a
SET farm_id = c.farm_id
FROM production_growing_cycles c
WHERE c.id = a.cycle_id AND a.farm_id IS NULL;

UPDATE production_plantings p
SET farm_id = c.farm_id
FROM production_growing_cycles c
WHERE c.id = p.cycle_id AND p.farm_id IS NULL;

UPDATE production_harvest_batch h
SET farm_id = c.farm_id
FROM production_growing_cycles c
WHERE c.id = h.cycle_id AND h.farm_id IS NULL;

ALTER TABLE production_allocations   ALTER COLUMN farm_id SET NOT NULL;
ALTER TABLE production_plantings     ALTER COLUMN farm_id SET NOT NULL;
ALTER TABLE production_harvest_batch ALTER COLUMN farm_id SET NOT NULL;

CREATE INDEX idx_production_allocations_farm_id   ON production_allocations (farm_id);
CREATE INDEX idx_production_plantings_farm_id     ON production_plantings (farm_id);
CREATE INDEX idx_production_harvest_batch_farm_id ON production_harvest_batch (farm_id);