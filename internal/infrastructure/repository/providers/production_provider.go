package providers

import (
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	domain "github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	inmemory "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/inmemory/production"
	postgres "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/postgres/production"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

func (p *productionProvider) ProviderName() string {
	return "production"
}

func NewProductionProvider(db uow.DB) repository.RepositoryProvider {
	return &productionProvider{db: db, inMemory: false}
}

type productionProvider struct {
	db          uow.DB
	inMemory    bool
	cycles      domain.GrowingCycleRepository
	allocations domain.AllocationRepository
	planting    domain.PlantingRepository
	harvest     domain.HarvestBatchRepository
}

func (p *productionProvider) Harvests() domain.HarvestBatchRepository {
	if p.harvest != nil {
		return p.harvest
	}
	if p.inMemory {
		p.harvest = inmemory.NewHarvestBatchRepository()
	} else {
		p.harvest = postgres.NewHarvestRepository(p.db)
	}
	return p.harvest
}

func (p *productionProvider) Planting() domain.PlantingRepository {
	if p.planting != nil {
		return p.planting
	}
	if p.inMemory {
		p.planting = inmemory.NewPlantingRepository()
	} else {
		p.planting = postgres.NewPlantingRepository(p.db)
	}
	return p.planting
}

func (p *productionProvider) Allocation() domain.AllocationRepository {
	if p.allocations != nil {
		return p.allocations
	}
	if p.inMemory {
		p.allocations = inmemory.NewAllocationRepository()
	} else {
		p.allocations = postgres.NewAllocationRepository(p.db)
	}
	return p.allocations
}

func (p *productionProvider) GrowingCycles() domain.GrowingCycleRepository {
	if p.cycles != nil {
		return p.cycles
	}
	if p.inMemory {
		p.cycles = inmemory.NewGrowingCycleRepository()
	} else {
		p.cycles = postgres.NewGrowingCycleRepository(p.db)
	}
	return p.cycles
}
