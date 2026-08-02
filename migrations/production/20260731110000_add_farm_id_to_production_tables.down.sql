ALTER TABLE production_allocations   DROP COLUMN IF EXISTS farm_id;
ALTER TABLE production_plantings     DROP COLUMN IF EXISTS farm_id;
ALTER TABLE production_harvest_batch DROP COLUMN IF EXISTS farm_id;