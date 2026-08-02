package production

import (
	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/allocation"
	growingcycle "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/growing_cycle"
	harvestbatch "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/harvest_batch"
	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/planting"
)

func scanGrowingCycle(root *growingcycle.GrowingCycle) []any {
	return []any{
		&root.ID,
		&root.FarmID,
		&root.CropID,
		&root.VarietyID,
		&root.ProtocolID,
		&root.Name,
		&root.Code,
		&root.Method,
		&root.Status,
		&root.Stage,
		&root.CreatedAt,
		&root.UpdatedAt,
	}
}

func scanAllocation(root *allocation.Allocation) []any {
	return []any{
		&root.ID,
		&root.FarmID,
		&root.CycleID,
		&root.ProductionUnitID,
		&root.Area,
		&root.StartedAt,
		&root.EndedAt,
		&root.CreatedAt,
		&root.UpdatedAt,
	}
}

func scanPlanting(root *planting.Planting) []any {
	return []any{
		&root.ID,
		&root.FarmID,
		&root.CycleID,
		&root.PlantedAt,
		&root.Quantity,
		&root.CreatedAt,
		&root.UpdatedAt,
	}
}

func scanHarvest(root *harvestbatch.HarvestBatch) []any {
	return []any{
		&root.ID,
		&root.FarmID,
		&root.CycleID,
		&root.HarvestedAt,
		&root.Quantity,
		&root.CreatedAt,
		&root.UpdatedAt,
	}
}
