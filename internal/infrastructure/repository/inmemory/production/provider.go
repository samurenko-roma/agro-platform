package production

import (
	"github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
)

type provider struct {
	cycles      repository.GrowingCycleRepository
	allocations repository.AllocationRepository
	plantings   repository.PlantingRepository
	harvests    repository.HarvestBatchRepository
}

func NewProvider() repository.ProductionProvider {
	return &provider{
		cycles:      NewGrowingCycleRepository(),
		allocations: NewAllocationRepository(),
		plantings:   NewPlantingRepository(),
		harvests:    NewHarvestBatchRepository(),
	}
}

func (p *provider) GrowingCycles() repository.GrowingCycleRepository { return p.cycles }
func (p *provider) Allocation() repository.AllocationRepository      { return p.allocations }
func (p *provider) Planting() repository.PlantingRepository          { return p.plantings }
func (p *provider) Harvests() repository.HarvestBatchRepository      { return p.harvests }
func (p *provider) ProviderName() string                             { return "production_inmemory" }

var _ repository.ProductionProvider = (*provider)(nil)
